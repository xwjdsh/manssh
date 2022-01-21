<template>
  <div @click="editRecord" :class="[edit ? 'edit' : '', 'record']">
    <h3>
      {{ record.alias }}
      <i @click="$emit('delete-record', record.alias)" class="fas fa-times"></i>
    </h3>
    <p>{{ record.day }}</p>
  </div>
</template>

<script>
export default {
  name: "Record",
  data() {
    return {
      edit: false,
    };
  },
  props: {
    record: Object,
  },
  methods: {
    editRecord() {
      this.edit = !this.edit;
      this.$emit("edit-record");
    },
  },
};
</script>

<style scoped>
.fas {
  color: red;
}
.record {
  background: #f4f4f4;
  margin: 5px;
  padding: 10px 20px;
  cursor: pointer;
}
.edit {
  border-left: 10px solid green;
}
.record h3 {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>