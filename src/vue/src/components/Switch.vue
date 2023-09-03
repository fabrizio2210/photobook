<script>
export default {
  props: {
    checkboxId: String,
    checked: Boolean,
  },
  methods: {
    changeCheck() {
      this.$emit("changeCheck", !this.checked);
    },
  },
  data() {return {};}
};
</script>
<template>
  <div class="can-toggle">
    <input
      :id="checkboxId"
      @click="changeCheck()"
      v-model="this.checked"
      type="checkbox"
    />
    <label :for="checkboxId">
      <div
        class="can-toggle__switch"
        data-checked="Sì"
        data-unchecked="No"
      ></div>
    </label>
  </div>
</template>

<style scoped>
body {
  padding: 30px;
}

.can-toggle {
  position: relative;
}
.can-toggle *,
.can-toggle *:before,
.can-toggle *:after {
  box-sizing: border-box;
}
.can-toggle input[type="checkbox"] {
  opacity: 0;
  position: absolute;
  top: 0;
  left: 0;
}
.can-toggle input[type="checkbox"][disabled] ~ label {
  pointer-events: none;
}
.can-toggle input[type="checkbox"][disabled] ~ label .can-toggle__switch {
  opacity: 0.4;
}
.can-toggle input[type="checkbox"]:checked ~ label .can-toggle__switch:before {
  content: attr(data-unchecked);
  left: 0;
}
.can-toggle input[type="checkbox"]:checked ~ label .can-toggle__switch:after {
  content: attr(data-checked);
}
.can-toggle label {
  user-select: none;
  position: relative;
  display: flex;
  align-items: center;
}
.can-toggle label .can-toggle__label-text {
  flex: 1;
  padding-left: 32px;
}
.can-toggle label .can-toggle__switch {
  position: relative;
}
.can-toggle label .can-toggle__switch:before {
  content: attr(data-checked);
  position: absolute;
  top: 0;
  text-transform: uppercase;
  text-align: center;
}
.can-toggle label .can-toggle__switch:after {
  content: attr(data-unchecked);
  position: absolute;
  z-index: 5;
  text-transform: uppercase;
  text-align: center;
  background: white;
  transform: translate3d(0, 0, 0);
}
.can-toggle input[type="checkbox"][disabled] ~ label {
  color: rgba(234, 0, 48, 0.1);
}
.can-toggle input[type="checkbox"]:focus ~ label .can-toggle__switch,
.can-toggle input[type="checkbox"]:hover ~ label .can-toggle__switch {
  background-color: #acb1bc; /*#ea0030;*/
}
.can-toggle input[type="checkbox"]:focus ~ label .can-toggle__switch:after,
.can-toggle input[type="checkbox"]:hover ~ label .can-toggle__switch:after {
  color: #7b5267; /* #b70026;*/
}
.can-toggle input[type="checkbox"]:hover ~ label {
  color: #8c6f87; /*#d1002b;*/
}
.can-toggle input[type="checkbox"]:checked ~ label:hover {
  color: #55bc49;
}
.can-toggle input[type="checkbox"]:checked:focus ~ label .can-toggle__switch,
.can-toggle input[type="checkbox"]:checked ~ label .can-toggle__switch {
  background-color: #b7ffa8; /*#70c767; box color green */
}
.can-toggle
  input[type="checkbox"]:checked:focus
  ~ label
  .can-toggle__switch:after,
.can-toggle input[type="checkbox"]:checked ~ label .can-toggle__switch:after {
  color: #4fb743; /*color inside */
}
.can-toggle input[type="checkbox"]:checked:hover ~ label .can-toggle__switch {
  background-color: #75dcab; /* #5fc054; color green when hoover*/
}
.can-toggle
  input[type="checkbox"]:checked:hover
  ~ label
  .can-toggle__switch:after {
  color: #47a43d;
}
.can-toggle label .can-toggle__label-text {
  flex: 1;
}
.can-toggle label .can-toggle__switch {
  transition: background-color 0.3s cubic-bezier(0, 1, 0.5, 1);
  background: #acb1bc; /* #ff0538;*/
}
.can-toggle label .can-toggle__switch:before {
  color: rgba(255, 255, 255, 0.1);
}
.can-toggle label .can-toggle__switch:after {
  -webkit-transition: -webkit-transform 0.3s cubic-bezier(0, 1, 0.5, 1);
  transition: transform 0.3s cubic-bezier(0, 1, 0.5, 1);
  color: #7b5267; /*#ea0030;*/
}
.can-toggle input[type="checkbox"]:focus ~ label .can-toggle__switch:after,
.can-toggle input[type="checkbox"]:hover ~ label .can-toggle__switch:after {
  box-shadow: 0 3px 3px rgba(0, 0, 0, 0.4);
}
.can-toggle input[type="checkbox"]:checked ~ label .can-toggle__switch:after {
  transform: translate3d(65px, 0, 0);
}
.can-toggle
  input[type="checkbox"]:checked:focus
  ~ label
  .can-toggle__switch:after,
.can-toggle
  input[type="checkbox"]:checked:hover
  ~ label
  .can-toggle__switch:after {
  box-shadow: 0 3px 3px rgba(0, 0, 0, 0.4);
}
.can-toggle label {
  font-size: 14px;
}
.can-toggle label .can-toggle__switch {
  height: 36px;
  flex: 0 0 134px;
  border-radius: 4px;
}
.can-toggle label .can-toggle__switch:before {
  left: 67px;
  font-size: 12px;
  line-height: 36px;
  width: 67px;
  padding: 0 12px;
}
.can-toggle label .can-toggle__switch:after {
  top: 2px;
  left: 2px;
  border-radius: 2px;
  width: 65px;
  line-height: 32px;
  font-size: 12px;
}
.can-toggle label .can-toggle__switch:hover:after {
  box-shadow: 0 3px 3px rgba(0, 0, 0, 0.4);
}
.can-toggle.can-toggle--size-small
  input[type="checkbox"]:focus
  ~ label
  .can-toggle__switch:after,
.can-toggle.can-toggle--size-small
  input[type="checkbox"]:hover
  ~ label
  .can-toggle__switch:after {
  box-shadow: 0 2px 2px rgba(0, 0, 0, 0.4);
}
.can-toggle.can-toggle--size-small
  input[type="checkbox"]:checked
  ~ label
  .can-toggle__switch:after {
  transform: translate3d(44px, 0, 0);
}
.can-toggle.can-toggle--size-small
  input[type="checkbox"]:checked:focus
  ~ label
  .can-toggle__switch:after,
.can-toggle.can-toggle--size-small
  input[type="checkbox"]:checked:hover
  ~ label
  .can-toggle__switch:after {
  box-shadow: 0 2px 2px rgba(0, 0, 0, 0.4);
}
.can-toggle.can-toggle--size-small label {
  font-size: 13px;
}
.can-toggle.can-toggle--size-small label .can-toggle__switch {
  height: 28px;
  flex: 0 0 90px;
  border-radius: 2px;
}
.can-toggle.can-toggle--size-small label .can-toggle__switch:before {
  left: 45px;
  font-size: 10px;
  line-height: 28px;
  width: 45px;
  padding: 0 12px;
}
.can-toggle.can-toggle--size-small label .can-toggle__switch:after {
  top: 1px;
  left: 1px;
  border-radius: 1px;
  width: 44px;
  line-height: 26px;
  font-size: 10px;
}
.can-toggle.can-toggle--size-small label .can-toggle__switch:hover:after {
  box-shadow: 0 2px 2px rgba(0, 0, 0, 0.4);
}
</style>
