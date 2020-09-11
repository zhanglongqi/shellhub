package sshd

import (
	"C"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	sshserver "github.com/gliderlabs/ssh"
	"github.com/shellhub-io/shellhub/agent/pkg/osauth"
	"github.com/sirupsen/logrus"
)

type sshConn struct {
	net.Conn
	closeCallback func(string)
	ctx           sshserver.Context
}

func (c *sshConn) Close() error {
	c.closeCallback(c.ctx.SessionID())
	return c.Conn.Close()
}

type Server struct {
	sshd              *sshserver.Server
	cmds              map[string]*exec.Cmd
	Sessions          map[string]net.Conn
	deviceName        string
	mu                sync.Mutex
	keepAliveInterval int
}

func NewServer(privateKey string, keepAliveInterval int) *Server {
	s := &Server{
		cmds:              make(map[string]*exec.Cmd),
		Sessions:          make(map[string]net.Conn),
		keepAliveInterval: keepAliveInterval,
	}

	s.sshd = &sshserver.Server{
		PasswordHandler:  s.passwordHandler,
		PublicKeyHandler: s.publicKeyHandler,
		Handler:          s.sessionHandler,
		RequestHandlers:  sshserver.DefaultRequestHandlers,
		ChannelHandlers:  sshserver.DefaultChannelHandlers,
		ConnCallback: func(ctx sshserver.Context, conn net.Conn) net.Conn {
			closeCallback := func(id string) {
				s.mu.Lock()
				defer s.mu.Unlock()

				if v, ok := s.cmds[id]; ok {
					v.Process.Kill()
					delete(s.cmds, id)
				}
			}

			return &sshConn{conn, closeCallback, ctx}
		},
	}

	s.sshd.SetOption(sshserver.HostKeyFile(privateKey))

	return s
}

func (s *Server) ListenAndServe() error {
	return s.sshd.ListenAndServe()
}

func (s *Server) HandleConn(conn net.Conn) {
	s.sshd.HandleConn(conn)
}

func (s *Server) SetDeviceName(name string) {
	s.deviceName = name
}

func (s *Server) sessionHandler(session sshserver.Session) {
	sspty, winCh, isPty := session.Pty()

	log := logrus.WithFields(logrus.Fields{
		"user": session.User(),
		"pty":  isPty,
	})

	log.Info("New session request")

	go StartKeepAliveLoop(time.Second*time.Duration(s.keepAliveInterval), session)

	if isPty {
		scmd := newShellCmd(s, session.User(), sspty.Term)

		pts, err := startPty(scmd, session, winCh)
		if err != nil {
			logrus.Warn(err)
		}

		u := osauth.LookupUser(session.User())

		uid, _ := strconv.Atoi(u.UID)

		os.Chown(pts.Name(), uid, -1)

		logrus.WithFields(logrus.Fields{
			"user": session.User(),
			"pty": pts.Name(),
			"remoteaddr": session.RemoteAddr().String(),
			"localaddr": session.LocalAddr().String(),
			}).Info("Session started")

		s.mu.Lock()
		s.cmds[session.Context().Value(sshserver.ContextKeySessionID).(string)] = scmd
		s.mu.Unlock()

		if err := scmd.Wait(); err != nil {
			logrus.Warn(err)
		}

		logrus.WithFields(logrus.Fields{
			"user": session.User(),
			"pty": pts.Name(),
			"remoteaddr": session.RemoteAddr().String(),
			"localaddr": session.LocalAddr().String(),
			}).Info("Session ended")
	} else {
		u := osauth.LookupUser(session.User())
		cmd := newCmd(u, "", "", s.deviceName, session.Command()...)

		stdout, _ := cmd.StdoutPipe()
		stdin, _ := cmd.StdinPipe()

		logrus.WithFields(logrus.Fields{
			"user": session.User(),
			"remoteaddr": session.RemoteAddr().String(),
			"localaddr": session.LocalAddr().String(),
			"Raw command": session.RawCommand(),
			}).Info("Command started")

		cmd.Start()

		go func() {
			if _, err := io.Copy(stdin, session); err != nil {
				fmt.Println(err)
			}
		}()

		go func() {
			if _, err := io.Copy(session, stdout); err != nil {
				fmt.Println(err)
			}
		}()

		cmd.Wait()

		logrus.WithFields(logrus.Fields{
			"user": session.User(),
			"remoteaddr": session.RemoteAddr().String(),
			"localaddr": session.LocalAddr().String(),
			"Raw command": session.RawCommand(),
			}).Info("Command ended")
	}
}

func (s *Server) passwordHandler(ctx sshserver.Context, pass string) bool {
	log := logrus.WithFields(logrus.Fields{
		"user": ctx.User(),
	})

	ok := osauth.AuthUser(ctx.User(), pass)

	if ok {
		log.Info("Accepted password")
	} else {
		log.Info("Failed password")
	}

	return ok
}

func (s *Server) publicKeyHandler(_ sshserver.Context, _ sshserver.PublicKey) bool {
	return true
}

func (s *Server) CloseSession(id string) {
	if session, ok := s.Sessions[id]; ok {
		session.Close()
		delete(s.Sessions, id)
	}
}

func newShellCmd(s *Server, username, term string) *exec.Cmd {
	shell := os.Getenv("SHELL")

	u := osauth.LookupUser(username)

	if shell == "" {
		shell = u.Shell
	}

	if term == "" {
		term = "xterm"
	}

	cmd := newCmd(u, shell, term, s.deviceName, shell, "--login")

	return cmd
}
