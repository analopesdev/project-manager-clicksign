<template>
  <li class="skeleton-card" aria-hidden="true">
    <q-skeleton square class="cover" />

    <div class="body">
      <q-skeleton type="text" class="title" />
      <q-skeleton type="text" class="client" />

      <div class="dates">
        <div v-for="date in 2" :key="date" class="date">
          <q-skeleton square class="date-icon" />
          <q-skeleton type="text" class="date-text" />
        </div>
      </div>
    </div>
  </li>
</template>

<style scoped lang="scss">
.skeleton-card {
  overflow: hidden;
  background: $surface;
  border: 1px solid $border;
  border-radius: $radius-lg;
  list-style: none;
}

.cover {
  width: 100%;
  height: auto;
  aspect-ratio: 3 / 2;
}

.body {
  display: flex;
  flex-direction: column;
  gap: $space-2;
  padding: $space-5;
}

.title {
  width: 60%;
  height: 25px;
}

.client {
  width: 40%;
  height: 20px;
}

.dates {
  display: flex;
  flex-direction: column;
  gap: $space-3;
  margin-top: $space-2;
  padding-top: $space-4;
  border-top: 1px solid $divider;
}

.date {
  display: flex;
  align-items: center;
  gap: $space-4;
}

.date-icon {
  flex: none;
  width: $icon-lg;
  height: $icon-lg;
}

.date-text {
  width: 55%;
  height: 20px;
}
</style>
