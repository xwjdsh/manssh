<template>
  <div :key="record.id" v-for="record in records">
    <Record
      @delete-record="$emit('delete-record', record.alias)"
      @edit-record="$emit('edit-record', record.alias)"
      :record="record"
    />
  </div>
</template>


<script>
import Record from "./Record.vue";
export default {
  name: "Records",
  components: {
    Record,
  },
  props: {
    records: Array,
  },
  emits: ["delete-record", "edit-record"],
};
</script>