<!-- @format -->

<template>
  <div>
    <div
      v-if="isDragging"
      @dragenter="onHoverDetectorEnter(true)"
      @dragleave="onHoverDetectorLeave(true)"
      class="hover-detector top"
    ></div>
    <div
      v-if="isDragging"
      @dragenter="onHoverDetectorEnter(false)"
      @dragleave="onHoverDetectorLeave(false)"
      class="hover-detector bottom"
    ></div>

    <SearchBar v-if="search" class="searchbar" @close="search = false" @input="onSearchInput">
      <b-dropdown :text="`Sorted by: ${sortByText}`" class="my-auto mr-3">
        <b-dropdown-item @click="onSortBy('custom')">Custom</b-dropdown-item>
        <b-dropdown-item @click="onSortBy('created')">Created Date</b-dropdown-item>
        <b-dropdown-item @click="onSortBy('title')">Title</b-dropdown-item>
      </b-dropdown>
    </SearchBar>
    <InfoBubble ref="info" color="orange" @hides="onInfoClose">
      <p>
        Searching for a specific page? Press
        <b>CTRL + F</b>!
      </p>
    </InfoBubble>
    <div class="page-container" :style="{ 'padding-top': search ? '75px' : '0' }">
      <draggable
        :list="pages"
        :disabled="search"
        chosenClass="chosen"
        @start="isDragging = true"
        @end="isDragging = false"
        @update="onUpdate"
      >
        <Page
          v-for="p in pagesVisible"
          :key="p.uid"
          :uid="p.uid"
          :title="p.title"
          :champs="p.champions.join(' ')"
          :primary="p.primary.tree"
          :secondary="p.secondary.tree"
          :prows="p.primary.rows.join(' ')"
          :srows="p.secondary.rows.join(' ')"
          :perks="p.perks.rows.join(' ')"
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
      pages: [],
      pagesVisible: [],
      search: false,
      sortBy: 'created',

      isDragging: false,
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
    },
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
        this.search = true;
        event.preventDefault();
      }
    },

    onEscapePress(event) {
      if (event.keyCode == 27) {
        this.search = false;
        this.pagesVisible = this.pages;
        event.preventDefault();
      }
    },

    onSearchInput(e) {
      let txt = e.text.toLowerCase();
      if (txt.length === 0) {
        this.pagesVisible = this.pages;
        return;
      }
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

  created: function() {
    this.sortBy = this.$route.query.sortBy;

    if (!this.sortBy) {
      this.sortBy = window.localStorage.getItem('sort-pages-by') || 'created';
    }

    this.reload();

    Utils.setWindowListener('keydown', this.onSearchPress);
    Utils.setWindowListener('keydown', this.onEscapePress);
  },

  destroyed: function() {
    Utils.removeWindowListener('keydown', this.onSearchPress);
    Utils.removeWindowListener('keydown', this.onEscapePress);
  },

  mounted() {
    if (!window.localStorage['info-page-search']) {
      setTimeout(this.$refs.info.show, 3000);
    }
  },
};
</script>

<style scoped>
/** @format */

.page-container {
  transition: all 0.25s ease-in-out;
}

.champ-header {
  display: flex;
}

.champ-header > img {
  border-radius: 50%;
  margin-right: 15px;
}

.searchbar {
  position: fixed;
  left: 20px;
  right: 20px;
  z-index: 5;
  transition: all 0.25s ease-in-out;
}
</style>
