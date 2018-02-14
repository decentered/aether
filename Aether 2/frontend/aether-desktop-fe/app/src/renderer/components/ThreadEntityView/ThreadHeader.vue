<template>
  <div class="thread-header">
    <div class="header-text">
      Gardait ces couvert melezes une. Avance certes fit mirent fronts jet plates. En ma pont tu tira au prit nous meme.
    </div>
    <div class="thread-body">
      Oie ras vif moi mort fin abat. Encore voyage le intime pleine oh eglise je mollir. Une bonne savez oui osait sapin vin pieds mines. Mais ses chez nul nous long. Lilas he fonds glace au halte temps la. Six fanatiques ere pressaient entrainait souffrance sur dur. Metres gronde jet visite ces trotte oui ses. Arrachait oh desespoir descendit pu chaclosah. Yeux elan sa vers doit pere la. Electrique on oh inattendus nationales singuliers convulsion maintenant.
      <br>
      Maintenant roc ecouterent descendons eclaireurs son remarquent. Egorger bambous la un dociles. Peu les reste jusqu ton gagne quels. Mais ifs chez moi sans dela. Est non ils galons astres clairs. Eu demain va le la allons blancs virent femmes adorer. Qu lilas petit voici haine abris en je nuits. Fuite monde au reste treve aides au. Vie soixante aisselle illumine six ces persuada preparer. Jet grave car non trois ronde mines les.
      <br>
      Net remporte lampions avancent ses mettions remettre tenacite. Le vont ah ii crie nord ca rues paix cree. Caracolent gourmettes ifs sur inattendus condamnait indulgence foi. Frenes baisse poemes net ils disait qui. Et arracha on ciments cartons et la fatigue falloir. Ici continue rizieres les feerique. Murmure ii presque le annonce du dansent. Meme epis jet tout prit dit.
      <br>
      Sol que dentelees moi caractere concierge age. Ravins ans venait propos bottes fer triste furent. Nul des sein fut asie robe eue. Ah miserable gachettes du bourreaux le. Ai je ramassa groupes extreme la dociles. Depenser arbustes poternes fumantes blanches sa or la ou actrices.
    </div>
    <div class="meta">
      <Tickers></Tickers>
      <Byline></Byline>
      <Actions></Actions>
      <Datetime></Datetime>
    </div>
  </div>

</template>

<script>
  import Tickers from './../Shared/Tickers'
  import Byline from './../Shared/Byline'
  import Actions from './../Shared/Actions'
  import Datetime from './../Shared/Datetime'
  export default {
    components: {
      Tickers,
      Byline,
      Actions,
      Datetime
    },
    data () {
      return {}
    },
    methods: {}
  }
</script>

<style scoped lang="scss">
  .thread-header
  {
    border-bottom:1px solid #ccc;
    .header-text
    {
      font-size:24px;
      font-weight:500;
    }
    .thread-body
    {
      padding-top:12px;
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
