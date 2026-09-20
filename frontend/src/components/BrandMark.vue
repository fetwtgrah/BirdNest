<script setup lang="ts">
import { Feather } from 'lucide-vue-next'
</script>

<template>
  <RouterLink class="brand" to="/" aria-label="BirdNest 首页">
    <span class="brand__mark">
      <Feather :size="19" :stroke-width="2" />
    </span>
    <span>BirdNest</span>
  </RouterLink>
</template>
