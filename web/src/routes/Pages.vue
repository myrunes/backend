<!-- @format -->

<template>
  <div>
    <div
      v-if="isDragging"
      class="hover-detector top"
      @dragenter="onHoverDetectorEnter(true)"
      @dragleave="onHoverDetectorLeave(true)"
    ></div>
    <div
      v-if="isDragging"
      class="hover-detector bottom"
      @dragenter="onHoverDetectorEnter(false)"
      @dragleave="onHoverDetectorLeave(false)"
    ></div>

    <InfoBubble ref="info" color="orange" @hides="onInfoClose">
      <p>
        Searching for a specific page? Press
        <b>CTRL + F</b>!
      </p>
    </InfoBubble>

    <div class="d-flex mb-2">
      <SearchBar
        class="searchbar mb-2"
        ref="searchBar"
        placeholder="Search for page name or champion"
        @input="onSearchInput"
      />
      <b-dropdown :text="`Sorted by: ${sortByText}`" class="drop-down" toggle-class="drop-down-btn">
        <b-dropdown-item @click="onSortBy('custom')">Custom</b-dropdown-item>
        <b-dropdown-item @click="onSortBy('created')">Created Date</b-dropdown-item>
        <b-dropdown-item @click="onSortBy('title')">Title</b-dropdown-item>
      </b-dropdown>
    </div>

    <div class="page-container">
      <h3
        v-if="pages !== null && pages.length < 1"
        class="no-pages"
      >You have not created any pages yet. : (</h3>

      <draggable
        :list="pages"
        :disabled="isSearch"
        chosen-class="chosen"
        @start="isDragging = true"
        @end="isDragging = false"
        @update="onUpdate"
      >
        <Page
          v-for="p in pagesVisible"
          :key="p.uid"
          :uid="p.uid"
          :title="p.title"
          :champs="p.champions"
          :primary="p.primary.tree"
          :secondary="p.secondary.tree"
          :prows="p.primary.rows"
          :srows="p.secondary.rows"
          :perks="p.perks.rows"
          @delete="deleted"
        />
      </draggable>
    </div>

    <div class="ctrl-btns">
      <button
        class="btn-slide btn-new"
        @click="$router.push({ name: 'RunePage', params: { uid: 'new' } })"
      >+</button>
    </div>
  </div>
</template>

<script>
/** @format */

import Rest from '../js/rest';
import Utils from '../js/utils';
import Page from '../components/Page';
import SearchBar from '../components/SearchBar';
import InfoBubble from '../components/InfoBubble';

import Draggable from 'vuedraggable';

export default {
  name: 'Champ',

  components: {
    Page,
    SearchBar,
    InfoBubble,
    Draggable,
  },

  data: function() {
    return {
      pages: null,
      pagesVisible: [],
      sortBy: 'created',

      isDragging: false,
      isSearch: false,
      scrollTimer: null,
    };
  },

  computed: {
    sortByText: function() {
      switch (this.sortBy) {
        case 'created':
          return 'Created Date';
        case 'custom':
          return 'Custom';
        case 'title':
          return 'Title';
      }

      return 'Default';
    },
  },

  created: function() {
    this.sortBy = this.$route.query.sortBy;

    if (!this.sortBy) {
      this.sortBy = window.localStorage.getItem('sort-pages-by') || 'created';
    }

    this.reload();

    Utils.setWindowListener('keydown', this.onSearchPress);
  },

  destroyed: function() {
    Utils.removeWindowListener('keydown', this.onSearchPress);
  },

  mounted() {
    if (!window.localStorage['info-page-search']) {
      setTimeout(this.$refs.info.show, 3000);
    }
  },

  methods: {
    reload() {
      Rest.getPages(this.sortBy)
        .then((res) => {
          if (!res.body) return;
          this.pages = this.pagesVisible = res.body.data;
        })
        .catch(console.error);
    },

    onUpdate(e) {
      this.sortBy = 'custom';
      window.localStorage.setItem('sort-pages-by', this.sortBy);
      console.log(this.pages.map((p) => p.uid));
      Rest.setPageOrder(this.pages.map((p) => p.uid));
    },

    deleted() {
      this.reload();
    },

    onSearchPress(event) {
      if (event.keyCode == 114 || (event.ctrlKey && event.keyCode == 70)) {
        this.$refs.searchBar.focus();
        event.preventDefault();
      }
    },

    onSearchInput(e) {
      let txt = e.text.toLowerCase();
      if (txt.length === 0) {
        this.pagesVisible = this.pages;
        this.isSearch = false;
        return;
      }
      this.isSearch = true;
      this.pagesVisible = this.pages.filter((p) => {
        return (
          p.title.toLowerCase().includes(txt) ||
          p.champions.find((c) => c.includes(txt))
        );
      });
    },

    onInfoClose(selfTriggered) {
      if (selfTriggered) {
        window.localStorage['info-page-search'] = '1';
      }
    },

    onSortBy(sortBy) {
      this.sortBy = sortBy;
      this.reload();
      window.localStorage.setItem('sort-pages-by', sortBy);
    },

    onHoverDetectorEnter(isTop) {
      if (!this.isDragging) return;

      if (isTop) {
        this.scrollTimer = setInterval(() => {
          window.scrollBy({
            top: -100,
          });
        }, 100);
      } else {
        this.scrollTimer = setInterval(() => {
          window.scrollBy({
            top: 100,
          });
        }, 100);
      }
    },

    onHoverDetectorLeave(isTop) {
      clearInterval(this.scrollTimer);
    },
  },
};
</script>

<style scoped>
/** @format */

.champ-header {
  display: flex;
}

.champ-header > img {
  border-radius: 50%;
  margin-right: 15px;
}

.no-pages {
  font-style: italic;
  text-align: center;
  margin-top: 30vh;
}

.drop-down {
  margin: auto 0 auto 10px;
}
</style>