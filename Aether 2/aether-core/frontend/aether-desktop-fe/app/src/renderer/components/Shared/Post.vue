<template>
  <div class="post" :class="'depth-'+depth">
    <div class="post-body">
    Net remporte lampions avancent ses mettions remettre tenacite. Le vont ah ii crie nord ca rues paix cree. Caracolent gourmettes ifs sur inattendus condamnait indulgence foi. Frenes baisse poemes net ils disait qui. Et arracha on ciments cartons et la fatigue falloir. Ici continue rizieres les feerique. Murmure ii presque le annonce du dansent. Meme epis jet tout prit dit.
    <br>
    Sol que dentelees moi caractere concierge age. Ravins ans venait propos bottes fer triste furent. Nul des sein fut asie robe eue. Ah miserable gachettes du bourreaux le. Ai je ramassa groupes extreme la dociles. Depenser arbustes poternes fumantes blanches sa or la ou actrices.</div>
    <div class="meta">
      <Tickers></Tickers>
      <Byline></Byline>
      <Actions></Actions>
    </div>
  </div>
</template>

<script>
  import Tickers from './../Shared/Tickers'
  import Byline from './../Shared/Byline'
  import Actions from './../Shared/Actions'
  export default {
    components: {
      Tickers,
      Byline,
      Actions
    },
    data () {
      return {}
    },
    methods: {},
    props: ['depth']
  }
</script>

<style scoped lang="scss">
  $deepest-indent:150px;
  .post
  {
    margin-left: $deepest-indent;
    /*This is important to have because*/
    /*if it's too deep, it'll stay at the*/
    /*max and not go back to zero indent.*/
    &.depth-0
    {
      margin-left:$deepest-indent/8*0;
    }
    &.depth-1
    {
      margin-left:$deepest-indent/8*1;
    }
    &.depth-2
    {
      margin-left:$deepest-indent/8*2;
    }
    border-bottom:1px solid #ccc;
    .post-body
    {
      margin-top:5px;
      margin-bottom:5px;
      padding-top:8px;
      line-height:140%;
      br
      {
        display: block;
        content: "";
        margin-top: 8px;
      }
    }
    .meta
    {
      display:flex;
      line-height:36px;
      & > * /* All direct descentants */
      {
        margin-right:20px;
      }
    }
  }
</style>
