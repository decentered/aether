<template>
  <div class="settings-sublocation">
    <a-markdown :content="headline"></a-markdown>
    <a-markdown :content="intro"></a-markdown>
    <a-markdown :content="content"></a-markdown>
  </div>
</template>

<script lang="ts">
  export default {
    name: 'advancedsettings',
    data() {
      return {
        headline: headline,
        intro: intro,
        content: content,
      }
    }
  }
  // These are var's and not let's because lets are defined only from the point they're in the code, and vars are defined for the whole scope regardless of where they are.
  var headline = '# Advanced Settings'
  var intro =
    `**This describes the knobs in more detail and provides instructions on how to change them.**

These descriptions are primarily intended for power-users, and all come with sane defaults. If none of this makes sense to you, you can happily ignore them.
  `
  var content =
    `


### Manually changing settings

This app attempts to provide sane defaults, so you should hopefully not need this. But in the case you are curious, you can find the config files in your application user directory, named backend and frontend configs, and make your changes there. Make sure that the app is fully shut down before editing.

Descriptions of the more important settings are below. For the rest, the descriptions can be found [here](https://github.com/nehbit/aether/blob/master/Aether%202/services/configstore/permanent.go).

Heads up, there are some settings, if misconfigured, can get your machine and user key permanently banned by other nodes in the network. As a general rule, the settings that modify local behaviour are probably safe to fiddle with, while ones that relate to network behaviour are officially The Danger Zone™.

### Defaults

| Knob | Description | Value |
|--- | --- |--- |
| Maximum disk space for database | The maximum disk space the backend database is allowed to take. Whenever the app reaches this threshold, it starts to delete from history to remain under the threshold. <br><br> Mind that this is not the entire disk space used by this app. The other things using the cache are the pre-baked HTTP caches that are used to ease outbound serves, and the frontend key-value store that holds precompiled graph of human-readable objects, such as boards, threads, users. <br><br> Both of these take an additional 10-15% of the database size, so at the maximum DB size of 10 Gb, the total disk use would be something around 13 Gb. | 10 Gb
| Local memory | For how long the local node remembers data for, in absence of disk space pressure. <br><br> This means the data will be deleted at the 6-month mark even if the maximum disk space is not reached. If it is, the local memory will be less. In other words, you can have an 6-month local memory but give the maximum disk space a value of 1 Gb, you will likely have much less than 6 months worth of content.  <br><br> This is somewhat akin to a reference-counting garbage collector, in that it will traverse the content graph, and if an object is older than 6 months but there exists a vertex that links to a newer-than-6-months graph node, it will not be deleted.  <br><br> An example of this is board graph nodes, which will not be deleted as long as someone has posted a thread in them in the last 6 months. | 6 months |
| Neighbourhood size | This is the size of the local node's neighbourhood that it will try to keep in sync with. <br><br> This is a push/pop stack, and at the end of a cycle, the oldest neighbour will be evicted, and a new one that wasn't synced before will be added to it. This means the local node will be connecting to known nodes for the 90% of the time, and new ones 10% of the time. | 10
| Tick duration | This is the base unit of time. The local node will attempt to establish a connection to a standard live node every tick by popping a candidate from the neighbourhood.  <br><br> Every 10 ticks the neighbourhood is cycled by one (see 'Neighbourhood size' above), every 60 ticks the static nodes will be hit, and every 360 ticks, bootstrap nodes. | 60 seconds
| Reverse open | Reverse opens are a way to 'request' a connection from a remote node by connecting to that remote and passing over a raw Mim request asking for a reverse-connect using the same TCP socket.  <br><br> This is useful for nodes that are behind firewalls and uncooperating NATs. Without this, no other nodes would be able to connect to them directly in another way because of the firewall, rendering the content they create unable to reach the network.  <br><br> In case of erratic behaviour, this is a good first thing to consider disabling (as long as your network is configured right, and UPNP can port-map your router). | Enabled
| Maximum address table size | How many other nodes' addresses will be kept in the database. <br><br> Whenever this threshold is crossed, the addresses with the oldest last successful connection timestamp will be purged from memory. | 1000
| Maximum simultaneous inbound connections | The number of remotes that can be syncing with the local node at the same time. <br><br> Mind that the inbound and outbounds are different types of syncs, because syncs are one-way pulls. A node syncing with you doesn't mean that you get the changes on that node, it just means the remote node gets the changes in yours. This improves security, since no one can 'push' data into your machine. A result of this is that it is imperative that other nodes are able to connect to you, because if they do not, the content you create will never be able to leave your node, and reach the network. <br><br> Therefore setting this value to 0 might at first seem like a good hack to reduce bandwidth use at the expense of others, but it will also render you effectively invisible to everyone. The larger this value is, the faster your content will reach other users of the network, up to the point that your uplink bandwidth, CPU or disk is saturated to such a degree that the remotes are abandoning syncs with you because it takes so long. <br><br> If the app is taxing your computer too much, this is a good value to try reducing one by one. The default value is chosen as a balance between network connectivity and system resource use, and assuming you have a CPU made in the last decade, it should not be taxing your CPU very much, if at all. | 5
| Maximum simultaneous outbound connections | The number of remotes that the local node be syncing with at the same time. | 1



  `
</script>

<style lang="scss" scoped>
  @import "../../../scss/globals";
  .settings-sublocation {
    color: $a-grey-600;
    .markdowned {
      &:first-of-type {
        margin-bottom: 0;
      }
      margin-bottom: 40px;
    }
  }
</style>