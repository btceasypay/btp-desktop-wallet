<template>
  <el-container>
    <el-header class="cheader">
      <el-row type="flex">
        <el-col :span="6">
          <h2>首页</h2>
        </el-col>
        <el-col :span="12"></el-col>
        <el-col :span="6"></el-col>
      </el-row>
    </el-header>
    <el-main class="cmain">
      <div class="welcome">
        <img src="../assets/logo.png" />
        <h1>欢迎使用 Btp Wallet</h1>

        <div>
          <el-row type="flex">
            <el-col :span="4">
              <div class="rr">Source Code:</div>
            </el-col>
            <el-col :span="16">
              <a
                href="https://github.com/btceasypay/btp-desktop-wallet"
              >https://github.com/btceasypay/btp-desktop-wallet</a>
            </el-col>
          </el-row>
          <el-row type="flex">
            <el-col :span="4">
              <div class="rr">Issues::</div>
            </el-col>
            <el-col :span="16">
              <a
                href="https://github.com/btceasypay/btp-desktop-wallet/issues"
              >https://github.com/btceasypay/btp-desktop-wallet/issues</a>
            </el-col>
          </el-row>
          <el-row type="flex">
            <el-col :span="4">
              <div class="rr">Btp:</div>
            </el-col>
            <el-col :span="16">
              <a href="https://github.com/Btp">https://github.com/Btp</a>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-main>
  </el-container>
</template>
<style>
</style>

<script>
export default {
  data() {
    return {};
  },
  mounted() {
    this.checkWalletStats();
  },
  methods: {
    openWallet() {
      let _this = this;
      _this.$emit("walletPasswordDlg", "ui_openWallet", result => {
        if (!result) {
          _this.openWallet();
          return;
        }
        _this.checkWalletStats();
      });
    },
    checkWalletStats() {
      let _this = this;
      _this.$emit("getWalletStats", stat => {
        if (stat == "nil") {
          _this.$emit("createWalletDlg");
          return;
        }
        if (stat == "closed") {
          _this.openWallet();
          return;
        }
        _this.$emit("walletOk");
      });
    }
  }
};
</script>

