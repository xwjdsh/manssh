<template>
  <div class="container">
    <Header
      @toggle-add-record="toggleAddRecord"
      title="SSH aliases"
      :showAddRecord="showAddRecord"
    />
    <div v-show="showAddRecord">
      <AddRecord @add-record="addRecord" />
    </div>
    <Records
      @delete-record="deleteRecord"
      @edit-record="editRecord"
      :records="records"
    />
  </div>
</template>

<script>
import Header from "./components/Header.vue";
import Records from "./components/Records.vue";
import AddRecord from "./components/AddRecord.vue";
export default {
  name: "App",
  data() {
    return {
      records: [],
      showAddRecord: false,
    };
  },
  components: {
    Header,
    Records,
    AddRecord,
  },
  methods: {
    toggleAddRecord() {
      this.showAddRecord = !this.showAddRecord;
    },
    addRecord(record) {
      this.records = [...this.records, record];
    },
    deleteRecord(alias) {
      if (confirm("Are you sure?")) {
        this.records = this.records.filter((record) => record.alias !== alias);
      }
    },
    editRecord(id) {
      console.log(111);
    },
    async fetchRecords() {
      const res = await fetch("api/records");
      const obj = await res.json();
      return obj.data;
    },
  },
  async created() {
    this.records = await this.fetchRecords();
  },
};
</script>

<style>
@import url("https://fonts.googleapis.com/css2?family=Poppins:wght@300;400&display=swap");
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
body {
  font-family: "Poppins", sans-serif;
}
.container {
  max-width: 800px;
  margin: 30px auto;
  overflow: auto;
  min-height: 300px;
  border: 1px solid steelblue;
  padding: 30px;
  border-radius: 5px;
}
.btn {
  display: inline-block;
  background: #000;
  color: #fff;
  border: none;
  padding: 10px 20px;
  margin: 5px;
  border-radius: 5px;
  cursor: pointer;
  text-decoration: none;
  font-size: 15px;
  font-family: inherit;
}
.btn:focus {
  outline: none;
}
.btn:active {
  transform: scale(0.98);
}
.btn-block {
  display: block;
  width: 100%;
}
</style>
