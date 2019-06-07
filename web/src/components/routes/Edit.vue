<template>
  <div>
    <Banner v-if="banner.visible" :type="banner.type" class="mb-3">{{ banner.content }}</Banner>

    <div>
      <div class="position-relative mb-3">
        <input 
          type="text" 
          class="tb tb-title" 
          v-model="page.title"
          placeholder="TITLE"
          @change="titleChange"
        />
        <span class="tb w-100"/>
      </div>
      <div class="position-relative">
        <input 
          type="text" 
          ref="tbChamps"
          class="tb tb-champions"
          placeholder="CHAMPIONS"
          @change="champChange"
        />
        <span class="tb w-100"/>
      </div>
    </div>
    
    <!-- TREE PICKER -->
    <div class="my-4 mx-2">
      <a v-for="tree in Object.keys(runes.primary)" 
        :key="`tree-${tree}`"
        class="mr-2 bordered"
        :class="{disabled: 
          !(page.primary.tree === tree || 
          page.secondary.tree === tree)
        }"
        :name="changes.trees"
        @click="treeClick(tree)"
      >
        <img :src="`/assets/rune-avis/${tree}.png`"/>
      </a>
    </div>

    <!-- TREES -->
    <div class="row mx-1 mb-4">
      <!-- PRIMARY TREE -->
      <div class="col bg mr-4">
        <h3>PRIMARY TREE</h3>
        <div v-for="(row, rowIndex) in runes.primary[page.primary.tree]" 
          :key="`row-${rowIndex}`"
          :name="changes.primary"
          class="mb-3"
        >
          <a v-for="rune in row" :key="`rune-${rune}`" 
            class="mr-2 bordered"
            :class="{disabled: page.primary.rows[rowIndex] !== rune}"
            @click="primaryClick(rowIndex, rune)"
          >
            <img 
              :src="`/assets/rune-avis/${page.primary.tree}/${rune}.png`"
              width="60" height="60"
            />
          </a>
        </div>
      </div>
      <!-- SECONDARY TREE -->
      <div class="col bg">
        <h3>SECONDARY TREE</h3>
        <div v-for="(row, rowIndex) in runes.secondary[page.secondary.tree]" 
          :key="`row-${rowIndex}`"
          :name="changes.secondary"
          class="mb-3"
        >
          <a v-for="rune in row" :key="`rune-${rune}`" 
            class="mr-2 bordered"
            :class="{disabled: !page.secondary.rows.includes(rune)}"
            @click="secondaryClick(rowIndex, rune)"
          >
            <img 
              :src="`/assets/rune-avis/${page.secondary.tree}/${rune}.png`"
              width="60" height="60"
            />
          </a>
        </div>
      </div>
    </div>

    <!-- PERKS -->
    <div class="bg ml-1">
      <h3>PERKS</h3>
      <div v-for="(row, index) in runes.perks" 
        :key="`perk-row-${index}`"
        class="mb-3"
        :name="changes.perks"
      >
        <a v-for="perk in row" 
          :key="`perk-${perk}`"
          class="mr-3 bordered"
          :class="{disabled: page.perks.rows[index] !== perk}"
          @click="perkClick(index, perk)"
        >
          <img 
            :src="`/assets/rune-avis/perks/${perk}.png`"
            width="40" height="40"
            class="perk big-perk"
            :class="`perk-${perk}`"
          />
        </a>
      </div>
    </div>

    <div class="ctrl-btns">
      <button 
        class="btn-slide mr-3 btn-cancel shadow"
        @click="$router.back()"
      >
        CANCEL
      </button>
      <button 
        class="btn-slide btn-save shadow"
        @click="save"
      >
        SAVE
      </button>
    </div>
  </div>
</template>

<script>
import Rest from '../../js/rest';
import Banner from '../Banner';

export default {
  name: 'Edit',

  components: {
    Banner,
  },

  data: function() {
    return {
      uid: null,

      champs: [],
      runes: {
        perks: [],
        primary: {},
        secondary: {},
      },

      page: {
        title: '',
        champions: [],
        primary: {
          rows: [],
        },
        secondary: {
          rows: [],
        },
        perks: {
          rows: [],
        }
      },

      banner: {
        visible: false,
        type: 'error',
        content: '',
      },

      changes: {
        trees: 0,
        primary: 0,
        secondary: 0,
        perks: 0,
      }
    }
  },

  methods: {
    getSecRow(rune) {
      let t = this.runes.secondary[this.page.secondary.tree];
      for (let i in t) {
        if (t[i].find((r) => r === rune))
          return i;
      }
      return -1;
    },

    titleChange(e) {
      if (e.target.value.length < 1) {
        this.banner = {
          visible: true,
          type: 'error',
          content: 'Title can not be empty!',
        }
      } else {
        this.banner.visible = false;
      }
    },

    champChange(e) {
      let val = this.$refs.tbChamps.value;
      if (val.length < 1) {
        this.banner = {
          visible: true,
          type: 'error',
          content: 'Champion list can not be empty!',
        }
        return;
      }

      let champs = val.split(',').map((v) => v.trim().toLowerCase());
      let invalid = [];
      for (let c of champs) {
        if (!this.champs.includes(c)) {
          invalid.push(c);
        }
      }

      if (invalid.length > 0) {
        this.banner = {
          visible: true,
          type: 'error',
          content: `Champion list contains invalid entries: ${invalid.join(', ')}`,
        }
        return;
      }
      
      this.banner.visible = false;
      this.page.champions = champs;
    },

    treeClick(tree) {
      this.changes.trees++;

      if (this.page.secondary.tree) {
        if (tree === this.page.secondary.tree) {
          this.page.secondary.tree = null;
          this.page.secondary.rows = [];
        }
        return;
      }

      if (tree === this.page.primary.tree) {
        this.page.primary.tree = null;
        this.page.primary.rows = [];
        return;
      }

      if (this.page.primary.tree) {
        this.page.secondary.tree = tree;
        this.page.secondary.rows = [];
        return;
      }

      this.page.primary.tree = tree;
      this.page.primary.rows = [];
    },

    primaryClick(rowIndex, rune) {
      this.page.primary.rows[rowIndex] = rune;
      this.changes.primary++;
    },

    secondaryClick(rowIndex, rune) {
      if (this.getSecRow(rune) === this.getSecRow(this.page.secondary.rows[0]))
        return;
      
      if (this.page.secondary.rows[0] && this.page.secondary.rows[1]) {
        this.page.secondary.rows[1] = this.page.secondary.rows[0];
        this.page.secondary.rows[0] = rune;
      } else if (this.page.secondary.rows[0]) {
        this.page.secondary.rows[1] = rune;
      } else {
        this.page.secondary.rows[0] = rune;
      }

      this.changes.secondary++;
    },

    perkClick(index, perk) {
      this.page.perks.rows[index] = perk;
      this.changes.perks++;
    },

    save() {
      var method;
      if (this.uid === 'new') {
        method = Rest.createPage(this.page);
      } else {
        method = Rest.updatePage(this.uid, this.page);
      }

      method.then((res) => {
        this.banner = {
          visible: true,
          type: 'success',
          content: 'Page saved!',
        }
        if (this.uid === 'new') {
          this.uid = res.body.uid;
          this.$router.replace({ name: 'RunePage', params: { uid: this.uid } });
        }
        window.scrollTo(0, 0);
        setTimeout(() => this.banner.visible = false, 10000);
      }).catch((err) => {
        this.banner = {
          visible: true,
          type: 'error',
          content: `Error: ${err.message ? err.message : err}`,
        }
        setTimeout(() => this.banner.visible = false, 10000);
        window.scrollTo(0, 0);
        console.error(err);
      })
    }
  },

  created: function() {
    this.uid = this.$route.params.uid;
    
    Rest.getChamps().then((res) => {
      if (!res.body || !res.body.data) return;
      this.champs = res.body.data;
    }).catch(console.error);

    Rest.getRunes().then((res) => {

      if (!res.body) return;
      this.runes = res.body;

      if (this.uid !== 'new') {
        Rest.getPage(this.uid).then((res) => {
          if (!res.body) return;
            this.page = res.body;
            this.$refs.tbChamps.value = this.page.champions.join(', ');
        }).catch(console.error);
      }
    }).catch(console.error);
  },

  updated: function() {
    if (this.$route.query && this.$route.query.champ) {
      let champ = this.$route.query.champ;
      this.page.champions.push(champ);
      this.$refs.tbChamps.value = champ;
    }
  }
}

</script>


<style scoped>

a:hover {
  cursor: pointer;
}

.tb-title {
  font-family: 'Montserrat', sans-serif;
  font-size: 30px;
  text-align: left;
  width: 100%;
}

.tb-champions {
  text-align: left;
  width: 100%;
}

.disabled {
  opacity: .6;
}

.disabled > img {
  border: none !important;
}

.bordered > img {
  border: solid #03A9F4 3px;
  border-radius: 50%;
}

.ctrl-btns {
  position: fixed;
  bottom: 30px;
  right: 30px;
  font-size: 18px;
}

</style>
