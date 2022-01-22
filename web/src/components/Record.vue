<template>
  <div @click="editRecord" :class="[edit ? 'edit' : '', 'record']">
    <h3>
      {{ record.alias }}
      {{ record.connection ? "-> " + record.connection : "" }}
      <i @click="$emit('delete-record', record.alias)" class="fas fa-times"></i>
    </h3>
    own_config:
    <div
      class="own_config"
      :key="key"
      v-for="key in Object.keys(record.own_config)"
    >
      <p>{{ key }} = {{ record.own_config[key] }}</p>
    </div>
    implicit_config:
    <div :key="key" v-for="key in Object.keys(record.implicit_config)">
      <p>{{ key }} = {{ record.implicit_config[key] }}</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "Record",
  data() {
    return {
      connection: "",
      edit: false,
    };
  },
  methods: {},
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

p {
  margin-left: 40px;
}

.own_config {
  color: green;
}
</style>