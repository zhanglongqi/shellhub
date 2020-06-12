<template>
  <div class="all-content">
    <div class="content-left">
      <v-navigation-drawer
        floating
        permanent
        absolute
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
    </div>
    <div class="content-right">
      <v-container
        class="pa-8"
        fluid
      >
        <router-view />
      </v-container>
    </div>
  </div>
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

<style lang="stylus" scoped>
  .content-left{
    float: left;
    width: 30%;
  }
  .content-right{
    float: right;
    width: 70%;
  }
</style>
