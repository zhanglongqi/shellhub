<template>
  <v-container>
    <v-row no-gutters>
      <v-col
        cols="auto"
      >
        <v-card
          class="pa-2"
          tile
          :elevation="0"
        >
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
        </v-card>
      </v-col>

      <v-col 
        cols="auto"
      >
        <v-card
          class="pa-2"
          tile
          :elevation="0"
        >
          <v-container
            class=""
            fluid
          >
            <router-view />
          </v-container>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
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
