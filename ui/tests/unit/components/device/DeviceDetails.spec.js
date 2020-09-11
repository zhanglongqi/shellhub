import Vuex from 'vuex';
import { shallowMount, createLocalVue } from '@vue/test-utils';
import DeviceDetails from '@/components/device/DeviceDetails';

describe('DeviceDetails', () => {
  let wrapper;

  const localVue = createLocalVue();
  localVue.use(Vuex);

  const store = new Vuex.Store({
    namespaced: true,
    state: {
      device: {
        uid: 'a582b47a42d',
        name: '39-5e-2a',
        identity: {
          mac: '00:00:00:00:00:00',
        },
        info: {
          id: 'arch',
          pretty_name: 'Linux Mint 19.3',
          version: '',
        },
        public_key: '----- PUBLIC KEY -----',
        tenant_id: '00000000',
        last_seen: '2020-05-20T18:58:53.276Z',
        online: false,
        namespace: 'user',
      },
    },
    getters: {
      'devices/get': (state) => state.device,
    },
    actions: {
      'devices/get': () => {
      },
    },
  });

  beforeEach(() => {
    wrapper = shallowMount(DeviceDetails, {
      store,
      localVue,
      stubs: ['fragment'],
      mocks: {
        $route: {
          params: {
            id: 'a582b47a42d',
          },
        },
      },
    });
  });

  it('Is a Vue instance', () => {
    expect(wrapper).toBeTruthy();
  });
  it('Renders the component', () => {
    expect(wrapper.html()).toMatchSnapshot();
  });
  it('Compare data and method with store', () => {
    expect(wrapper.vm.device.uid).toEqual('a582b47a42d');
    expect(wrapper.vm.device.identity.mac).toEqual('00:00:00:00:00:00');
    expect(wrapper.vm.device.info.pretty_name).toEqual('Linux Mint 19.3');
    expect(wrapper.vm.convertDate).toEqual('Wednesday, May 20th 2020, 6:58:53 pm');
  });
  it('Has a elements field and compare with store', () => {
    expect(wrapper.find('[data-test="deviceUid-field"]').text()).toEqual('a582b47a42d');
    expect(wrapper.find('[data-test="deviceMac-field"]').text()).toEqual('00:00:00:00:00:00');
    expect(wrapper.find('[data-test="devicePrettyName-field"]').text()).toEqual('Linux Mint 19.3');
    expect(wrapper.find('[data-test="deviceConvertDate-field"]').text()).toEqual('Wednesday, May 20th 2020, 6:58:53 pm');
  });
});
