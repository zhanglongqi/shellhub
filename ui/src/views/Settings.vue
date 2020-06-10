<template>
  <v-app>
    <v-navigation-drawer
      floating
      permanent
    >
      <v-list>
        <v-list-item
          v-for="item in items"
          :key="item.title"
          :to="item.path"
          two-line
        >
          <v-list-item-content>
            <v-list-item-title v-text="item.title" />
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-content>
      <v-container
        class="pa-8"
        fluid
      >
        <router-view />
      </v-container>
    </v-content>
  </v-app>
</template>

<script>

export default {
  name: 'Settings',

  data() {
    return {
      drawer: true,
      clipped: false,
      copySnack: false,
      items: [
        {
          title: 'Profile',
          path: '/settings/profile'
        },
      ],
    };
  },
};
</script>
