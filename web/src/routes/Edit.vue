<!-- @format -->

<template>
  <div>
    <b-modal id="modalShare" title="Share Page" @ok="shareOk">
      <p v-if="!share.uid">
        This page has not been shared yet. Create a share link below.
      </p>
      <div v-else>
        <p class="m-0">Share link:</p>
        <p class="bg-ident">{{ `${getWindowLocation()}/p/${share.ident}` }}</p>
        <p>
          Visited
          <b>{{ share.accesses }}</b> times.
        </p>
      </div>

      <h5 class="mt-4">Max Uses</h5>
      <i>
        Number of times the link can be accessed. Set to -1 to set this
        infinite.
      </i>
      <b-form-input
        v-model="share.maxaccesses"
        type="number"
        min="-1"
        value="0"
      ></b-form-input>

      <h5 class="mt-4">Expires</h5>
      <i>
        Time at which the link will expire. Leave empty to set to never expire.
      </i>
      <b-row>
        <b-col>
          <b-form-input
            ref="shareDate"
            v-model="share._expires.date"
            type="date"
          ></b-form-input>
        </b-col>
        <b-col>
          <b-form-input
            ref="shareTime"
            v-model="share._expires.time"
            type="time"
          ></b-form-input>
        </b-col>
      </b-row>

      <b-button
        v-if="share.uid"
        variant="danger"
        class="w-100 mt-3 text-white"
        @click="resetShare"
        >RESET SHARE</b-button
      >
    </b-modal>

    <Banner ref="banner" class="mb-3"></Banner>

    <div>
      <div class="position-relative mb-3">
        <input
          v-model="page.title"
          type="text"
          class="tb tb-title"
          placeholder="TITLE"
          @change="titleChange"
        />
        <span class="tb w-100" />
      </div>
      <TagsInput
        ref="tagChamps"
        :tags="champs"
        :formatter="champFormatter"
        :filter="champFilter"
        @change="champsChanged"
      />
    </div>

    <!-- TREE PICKER -->
    <div class="my-4 mx-2">
      <a
        v-for="tree in runes.trees"
        :key="`tree-${tree.uid}`"
        class="mr-2 bordered"
        :class="{
          disabled: !(
            page.primary.tree === tree.uid || page.secondary.tree === tree.uid
          ),
        }"
        :name="changes.trees"
        @click="treeClick(tree.uid)"
      >
        <img :src="`/assets/rune-avis/${tree.uid}.png`" />
      </a>
    </div>

    <!-- TREES -->
    <div class="row mx-1 mb-4">
      <!-- PRIMARY TREE -->
      <div class="col bg mr-4">
        <h3>PRIMARY TREE</h3>
        <div
          v-for="(row, rowIndex) in getPrimaryTree(page.primary.tree)"
          :key="`row-${rowIndex}`"
          :name="changes.primary"
          class="mb-3"
        >
          <a
            v-for="rune in row.runes"
            :key="`rune-${rune.uid}`"
            class="mr-2 bordered rune"
            :class="{ disabled: page.primary.rows[rowIndex] !== rune.uid }"
            @click="primaryClick(rowIndex, rune.uid)"
          >
            <img
              :id="`rune-${rune.uid}`"
              :src="`/assets/rune-avis/${page.primary.tree}/${rune.uid}.png`"
              width="60"
              height="60"
            />
            <b-tooltip :target="`rune-${rune.uid}`" delay="500">
              <div class="rune-tool-tip">
                <h5>{{ rune.name }}</h5>
                <p>
                  {{
                    formatRuneText(
                      showFullInfo ? rune.longDesc : rune.shortDesc
                    )
                  }}
                </p>
              </div>
            </b-tooltip>
          </a>
        </div>
      </div>

      <!-- SECONDARY TREE -->
      <div class="col bg">
        <h3>SECONDARY TREE</h3>
        <div
          v-for="(row, rowIndex) in getSecondaryTree(page.secondary.tree)"
          :key="`row-${rowIndex}`"
          :name="changes.secondary"
          class="mb-3"
        >
          <a
            v-for="rune in row.runes"
            :key="`rune-${rune.uid}`"
            class="mr-2 bordered"
            :class="{ disabled: !page.secondary.rows.includes(rune.uid) }"
            @click="secondaryClick(rowIndex, rune.uid)"
          >
            <img
              :id="`rune-${rune.uid}`"
              :src="`/assets/rune-avis/${page.secondary.tree}/${rune.uid}.png`"
              width="60"
              height="60"
            />
            <b-tooltip :target="`rune-${rune.uid}`" delay="500">
              <div class="rune-tool-tip">
                <h5>{{ rune.name }}</h5>
                <p>
                  {{
                    formatRuneText(
                      showFullInfo ? rune.longDesc : rune.shortDesc
                    )
                  }}
                </p>
              </div>
            </b-tooltip>
          </a>
        </div>
      </div>
    </div>

    <!-- PERKS -->
    <div class="bg ml-1">
      <h3>PERKS</h3>
      <div
        v-for="(row, index) in runes.perks"
        :key="`perk-row-${index}`"
        class="mb-3"
        :name="changes.perks"
      >
        <a
          v-for="perk in row"
          :key="`perk-${perk}`"
          class="mr-3 bordered"
          :class="{ disabled: page.perks.rows[index] !== perk }"
          @click="perkClick(index, perk)"
        >
          <img
            :src="`/assets/rune-avis/perks/${perk}.png`"
            width="40"
            height="40"
            class="perk big-perk"
            :class="`perk-${perk}`"
          />
        </a>
      </div>
    </div>

    <Slider
      ref="sliderFullInfo"
      class="mt-3 ml-1"
      v-model="showFullInfo"
      @input="onFullInfoSliderInput"
      >Show full rune info</Slider
    >

    <div class="ctrl-btns">
      <button v-if="created" class="btn-slide mr-3 shadow" @click="shareOpen">
        SHARE
      </button>
      <button class="btn-slide mr-3 btn-cancel shadow" @click="$router.back()">
        CANCEL
      </button>
      <button class="btn-slide btn-save shadow" @click="save">SAVE</button>
    </div>
  </div>
</template>

<script>
/** @format */

import Rest from '../js/rest';
import Banner from '../components/Banner';
import TagsInput from '../components/TagsInput';
import Slider from '../components/Slider';

export default {
  name: 'Edit',

  components: {
    Banner,
    TagsInput,
    Slider,
  },

  data: function() {
    return {
      uid: null,

      created: false,

      champs: [],
      runes: {
        trees: [],
        perks: [],
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
        },
      },

      share: {
        maxaccesses: -1,
        _expires: {},
      },

      changes: {
        trees: 0,
        primary: 0,
        secondary: 0,
        perks: 0,
      },

      wasUpdated: false,

      showFullInfo: false,
    };
  },

  created: function() {
    this.uid = this.$route.params.uid;

    Rest.getChamps()
      .then((res) => {
        this.champs = res.body.data;
      })
      .catch(console.error);

    Rest.getRunes()
      .then((res) => {
        if (!res.body) return;
        this.runes = res.body;

        if (this.uid !== 'new') {
          Rest.getPage(this.uid)
            .then((res) => {
              if (!res.body) return;
              this.created = true;
              this.page = res.body;
              this.page.champions
                .map((c) => this.champs.find((cd) => cd.uid === c))
                .forEach((c) => this.$refs.tagChamps.append(c));
            })
            .catch(console.error);
        }
      })
      .catch(console.error);
  },

  mounted() {
    if (window.localStorage.getItem('show-full-rune-info') === '1') {
      this.showFullInfo = true;
      this.$refs.sliderFullInfo.set(true);
    }
  },

  updated: function() {
    if (!this.wasUpdated && this.$route.query && this.$route.query.champ) {
      let champ = this.champs.find((c) => c.uid === this.$route.query.champ);
      this.page.champions.push(champ.uid);
      this.$refs.tagChamps.append(champ);
      this.wasUpdated = true;
    }
  },

  methods: {
    getSecRow(rune) {
      let t = this.runes.trees
        .find((t) => t.uid === this.page.secondary.tree)
        .slots.slice(1);
      for (let i in t) {
        if (t[i].runes.find((r) => r.uid === rune)) return i;
      }
      return -1;
    },

    titleChange(e) {
      if (e.target.value.length < 1) {
        this.$refs.banner.show('error', 'Title can not be empty!', 10000, true);
      } else {
        this.$refs.banner.hide();
      }
    },

    champsChanged(champs) {
      this.page.champions = champs.map((c) => c.uid);
    },

    getPrimaryTree(uid) {
      const tree = this.runes.trees.find((t) => t.uid === uid);
      if (!tree) return [];
      return tree.slots;
    },

    getSecondaryTree(uid) {
      const tree = this.runes.trees.find((t) => t.uid === uid);
      if (!tree) return [];
      return tree.slots.slice(1);
    },

    treeClick(tree) {
      this.changes.trees++;
      if (this.page.secondary.tree) {
        if (tree === this.page.secondary.tree) {
          this.page.secondary.tree = null;
          this.page.secondary.rows = [];
        } else if (tree !== this.page.primary.tree) {
          this.page.secondary.tree = tree;
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
      this.changes.secondary++;

      if (
        this.getSecRow(rune) === this.getSecRow(this.page.secondary.rows[0])
      ) {
        this.page.secondary.rows[0] = rune;
        return;
      }

      if (
        this.getSecRow(rune) === this.getSecRow(this.page.secondary.rows[1])
      ) {
        this.page.secondary.rows[1] = rune;
        return;
      }

      if (this.page.secondary.rows[0] && this.page.secondary.rows[1]) {
        this.page.secondary.rows[1] = this.page.secondary.rows[0];
        this.page.secondary.rows[0] = rune;
      } else if (this.page.secondary.rows[0]) {
        this.page.secondary.rows[1] = rune;
      } else {
        this.page.secondary.rows[0] = rune;
      }
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
      method
        .then((res) => {
          this.$refs.banner.show('success', 'Page saved!', 10000, true);
          if (this.uid === 'new') {
            this.uid = res.body.uid;
            this.$router.replace({
              name: 'RunePage',
              params: { uid: this.uid },
            });
          }
          this.created = true;
          window.scrollTo(0, 0);
        })
        .catch((err) => {
          this.$refs.banner.show(
            'error',
            `Error: ${err.message ? err.message : err}`,
            10000,
            true
          );
          window.scrollTo(0, 0);
          console.error(err);
        });
    },

    shareOpen() {
      Rest.getShare(this.uid)
        .then((res) => {
          this.share = res.body.share;
          this.share._expires = {};
          this.$bvModal.show('modalShare');
        })
        .catch((err) => {
          if (err.code !== 404) {
            console.error(err);
          } else {
            this.share = {
              maxaccesses: -1,
              _expires: {},
            };
            this.$bvModal.show('modalShare');
          }
        });
    },

    shareOk() {
      this.share.maxaccesses = !this.share.maxaccesses
        ? 0
        : parseInt(this.share.maxaccesses);

      let exp = new Date(
        this.share._expires.date + ' ' + this.share._expires.time
      );
      this.share.expires = exp.toString() === 'Invalid Date' ? null : exp;

      if (this.share.uid) {
        Rest.updateShare(this.share)
          .then(() => {
            this.$refs.banner.show('success', 'Share updated!', 10000, true);
          })
          .catch((err) => {
            this.$refs.banner.show(
              'error',
              `An error occured during saving share status: ${
                err.message ? err.message : err
              }`,
              10000,
              true
            );
          });
      } else {
        this.share.page = this.uid;
        Rest.createShare(this.share)
          .then((res) => {
            if (res.body) {
              this.$refs.banner.show(
                'success',
                `Share successfully created. Sharelink is: ${this.getWindowLocation()}/p/${
                  res.body.ident
                }`,
                null,
                true
              );
            }
          })
          .catch((err) => {
            this.$refs.banner.show(
              'error',
              `An error occured during saving share status: ${
                err.message ? err.message : err
              }`,
              10000,
              true
            );
          });
      }
    },

    resetShare() {
      Rest.deleteShare(this.share)
        .then(() => {
          this.$refs.banner.show(
            'success',
            'Now, this page is private again and share link will not work anymore.',
            10000,
            true
          );
        })
        .catch(console.error);
      this.$bvModal.hide('modalShare');
      this.share = {
        maxaccesses: -1,
        _expires: {},
      };
    },

    getWindowLocation() {
      return window.location.origin;
    },

    champFormatter(c) {
      return c.name;
    },

    champFilter(c, q) {
      return c.name.toLowerCase().includes(q.toLowerCase());
    },

    formatRuneText(txt) {
      return txt.replace(/<.*?>/gm, '');
    },

    onFullInfoSliderInput() {
      window.localStorage.setItem(
        'show-full-rune-info',
        this.showFullInfo ? '1' : '0'
      );
    },
  },
};
</script>

<style scoped>
/** @format */

a:hover {
  cursor: pointer;
}

.rune {
  position: relative;
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
  opacity: 0.6;
}

.disabled > img {
  border: none !important;
}

.bordered > img {
  border: solid #03a9f4 3px;
  border-radius: 50%;
}

.bg-ident {
  background-color: rgba(0, 0, 0, 0.1);
  width: fit-content;
  padding: 5px 10px;
  border-radius: 5px;
}

.rune-tool-tip {
  text-align: left;
}
</style>
