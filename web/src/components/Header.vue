<template>
  <header>
    <h1>{{ title }}</h1>
    <Button
      @btn-click="$emit('toggle-add-record')"
      :text="showAddRecord ? 'Close' : 'Add Record'"
      :color="showAddRecord ? 'red' : 'green'"
    />
  </header>
</template>


<script>
import Button from "./Button.vue";
export default {
  name: "Header",
  components: {
    Button,
  },
  props: {
    title: String,
    showAddRecord: Boolean,
  },
};
</script>


<style scoped>
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>