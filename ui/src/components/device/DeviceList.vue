<template>
  <fragment>
    <v-card-text class="pa-0">
      <v-data-table
        class="elevation-1"
        :headers="headers"
        :items="getListDevices"
        :items-per-page="10"
        :footer-props="{'items-per-page-options': [10, 25, 50, 100]}"
        :server-items-length="getNumberDevices"
        :options.sync="pagination"
      >
        <template v-slot:item.online="{ item }">
          <v-icon
            v-if="item.online"
            color="success"
          >
            check_circle
          </v-icon>
          <v-tooltip
            v-else
            bottom
          >
            <template #activator="{ on }">
              <v-icon v-on="on">
                check_circle
              </v-icon>
            </template>
            <span>last seen {{ item.last_seen | moment("from", "now") }}</span>
          </v-tooltip>
        </template>

        <template v-slot:item.hostname="{ item }">
          <router-link :to="{ name: 'detailsDevice', params: { id: item.uid } }">
            {{ item.name }}
          </router-link>
        </template>

        <template v-slot:item.info.pretty_name="{ item }">
          <DeviceIcon :icon-name="item.info.id" />
          {{ item.info.pretty_name }}
        </template>

        <template v-slot:item.namespace="{ item }">
          <v-chip class="list-itens">
            {{ address(item) }}
            <v-icon
              v-clipboard="() => address(item)"
              v-clipboard:success="showCopySnack"
              small
              right
              @click.stop
            >
              mdi-content-copy
            </v-icon>
          </v-chip>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-tooltip bottom>
            <template #activator="{ on }">
              <v-icon
                class="icons"
                v-on="on"
                @click="detailsDevice(item)"
              >
                info
              </v-icon>
            </template>
            <span>Details</span>
          </v-tooltip>

          <TerminalDialog
            v-if="item.online"
            :uid="item.uid"
          />

          <DeviceDelete
            :uid="item.uid"
            @update="refresh"
          />
        </template>
      </v-data-table>
    </v-card-text>
  </fragment>
</template>
<script>

import TerminalDialog from '@/components/terminal/TerminalDialog';
import DeviceIcon from '@/components/device/DeviceIcon';
import DeviceDelete from '@/components/device/DeviceDelete';
import formatOrdering from '@/components/device/Device';

export default {
  name: 'DeviceList',

  components: {
    TerminalDialog,
    DeviceIcon,
    DeviceDelete,
  },

  mixins: [formatOrdering],

  data() {
    return {
      hostname: window.location.hostname,
      pagination: {},
      headers: [
        {
          text: 'Online',
          value: 'online',
          align: 'center',
        },
        {
          text: 'Hostname',
          value: 'hostname',
          align: 'center',
        },
        {
          text: 'Operating System',
          value: 'info.pretty_name',
          align: 'center',
          sortable: false,
        },
        {
          text: 'SSHID',
          value: 'namespace',
          align: 'center',
          sortable: false,
        },
        {
          text: 'Actions',
          value: 'actions',
          align: 'center',
          sortable: false,
        },
      ],
    };
  },

  computed: {
    getListDevices() {
      return this.$store.getters['devices/list'];
    },

    getNumberDevices() {
      return this.$store.getters['devices/getNumberDevices'];
    },
  },

  watch: {
    pagination: {
      handler() {
        this.getDevices();
      },
      deep: true,
    },
  },

  methods: {
    async getDevices() {
      let sortStatusMap = {};

      sortStatusMap = this.formatSortObject(this.pagination.sortBy[0], this.pagination.sortDesc[0]);

      const data = {
        perPage: this.pagination.itemsPerPage,
        page: this.pagination.page,
        filter: this.$store.getters['devices/getFilter'],
        status: 'accepted',
        sortStatusField: sortStatusMap.field,
        sortStatusString: sortStatusMap.statusString,
      };

      try {
        await this.$store.dispatch('devices/fetch', data);
      } catch {
        this.$store.dispatch('modals/showSnackbarErrorLoading', this.$errors.deviceList);
      }
    },

    detailsDevice(value) {
      this.$router.push(`/device/${value.uid}`);
    },

    address(item) {
      return `${item.namespace}.${item.name}@${this.hostname}`;
    },

    copy(device) {
      this.$clipboard(device.uid);
    },

    showCopySnack() {
      this.$store.dispatch('modals/showSnackbarCopy', this.$copy.deviceSSHID);
    },

    refresh() {
      this.getDevices();
    },
  },
};

</script>

<style scoped>

.list-itens {
  font-family: monospace;
}

.icons{
  margin-right: 4px;
}

</style>
