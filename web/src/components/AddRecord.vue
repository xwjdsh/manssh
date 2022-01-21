<template>
  <form @submit="onSubmit" class="add-form">
    <div class="form-control">
      <label>Task</label>
      <input type="text" v-model="text" name="text" placeholder="Add Record" />
    </div>
    <div class="form-control">
      <label>Day & Time</label>
      <input
        type="text"
        v-model="day"
        name="day"
        placeholder="Add Day & Time"
      />
    </div>
    <input type="submit" value="Save" class="btn btn-block" />
  </form>
</template>


<script>
export default {
  name: "AddRecod",
  data() {
    return {
      text: "",
      day: "",
    };
  },
  methods: {
    onSubmit(e) {
      e.preventDefault();

      if (!this.text) {
        alert("Please input text");
      }

      const newRecord = {
        text: this.text,
        day: this.day,
      };

      this.$emit("add-record", newRecord);

      this.text = "";
      this.day = "";
    },
  },
};
</script>

<style scoped>
.add-form {
  margin-bottom: 40px;
}
.form-control {
  margin: 20px 0;
}
.form-control label {
  display: block;
}
.form-control input {
  width: 100%;
  height: 40px;
  margin: 5px;
  padding: 3px 7px;
  font-size: 17px;
}
.form-control-check {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.form-control-check label {
  flex: 1;
}
.form-control-check input {
  flex: 2;
  height: 20px;
}
</style>