<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <i class="material-icons large">cloud</i>
        <h4>S3 对象存储管理</h4>
        <p>请登录以访问控制台</p>
      </div>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="row">
          <div class="input-field col s12">
            <i class="material-icons prefix">account_circle</i>
            <input 
              id="username" 
              type="text" 
              v-model="username" 
              required
              :disabled="isLoading"
            >
            <label for="username">用户名</label>
          </div>
        </div>
        
        <div class="row">
          <div class="input-field col s12">
            <i class="material-icons prefix">lock</i>
            <input 
              id="password" 
              type="password" 
              v-model="password" 
              required
              :disabled="isLoading"
            >
            <label for="password">密码</label>
          </div>
        </div>
        
        <div class="row">
          <div class="col s12">
            <button 
              type="submit" 
              class="btn waves-effect waves-light col s12"
              :disabled="isLoading"
              style="background-color: #1976D2;"
            >
              <i class="material-icons left" v-if="!isLoading">login</i>
              <i class="material-icons left" v-else>hourglass_empty</i>
              {{ isLoading ? '登录中...' : '登录' }}
            </button>
          </div>
        </div>
      </form>
      
      <div class="login-footer">
        <p class="grey-text">默认账号: admin / admin</p>
      </div>
    </div>
  </div>
</template>

<script>
/* global M */
export default {
  name: 'Login',
  data() {
    return {
      username: '',
      password: '',
      isLoading: false
    }
  },
  methods: {
    async handleLogin() {
      if (!this.username || !this.password) {
        M.toast({ html: '请输入用户名和密码', classes: 'red' });
        return;
      }
      
      this.isLoading = true;
      
      try {
        // 简单的验证逻辑 - 你可以根据需要修改
        if (this.username === 'admin' && this.password === 'admin') {
          // 登录成功
          localStorage.setItem('isLoggedIn', 'true');
          localStorage.setItem('username', this.username);
          
          M.toast({ html: '登录成功！', classes: 'green' });
          
          // 延迟跳转，让用户看到成功提示
          setTimeout(() => {
            this.$router.push('/');
          }, 1000);
        } else {
          // 登录失败
          M.toast({ html: '用户名或密码错误', classes: 'red' });
        }
      } catch (error) {
        console.error('登录错误:', error);
        M.toast({ html: '登录失败，请重试', classes: 'red' });
      } finally {
        this.isLoading = false;
      }
    }
  },
  mounted() {
    // 初始化 Materialize 的输入框
    M.updateTextFields();
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header i {
  color: #1976D2;
  margin-bottom: 15px;
}

.login-header h4 {
  margin: 10px 0;
  color: #333;
}

.login-header p {
  color: #666;
  margin: 0;
}

.login-form {
  margin-bottom: 20px;
}

.login-footer {
  text-align: center;
  border-top: 1px solid #eee;
  padding-top: 20px;
}

.login-footer p {
  margin: 0;
  font-size: 14px;
}

.input-field input:focus + label {
  color: #1976D2 !important;
}

.input-field input:focus {
  border-bottom: 1px solid #1976D2 !important;
  box-shadow: 0 1px 0 0 #1976D2 !important;
}

.input-field .prefix.active {
  color: #1976D2 !important;
}
</style> 