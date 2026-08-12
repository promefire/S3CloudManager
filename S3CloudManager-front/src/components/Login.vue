<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <i class="material-icons large">cloud</i>
        <h4>S3 Storage</h4>
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
              class="btn waves-effect waves-light col s12 btn-login"
              :disabled="isLoading"
            >
              <i class="material-icons left" v-if="!isLoading">login</i>
              <i class="material-icons left" v-else>hourglass_empty</i>
              {{ isLoading ? '登录中...' : '登录' }}
            </button>
          </div>
        </div>
      </form>
      
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
        if (this.username === 'admin' && this.password === 'promefire123') {
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
  background: var(--color-bg);
  padding: var(--spacing-xl);
}

.login-card {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-lg);
  padding: var(--spacing-3xl);
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: var(--spacing-2xl);
}

.login-header i {
  color: var(--color-primary);
  margin-bottom: var(--spacing-md);
}

.login-header h4 {
  margin: var(--spacing-sm) 0;
  color: var(--color-text);
  font-size: var(--font-size-xl);
  font-weight: 600;
}

.login-header p {
  color: var(--color-text-secondary);
  margin: 0;
  font-size: var(--font-size-md);
}

.login-form {
  margin-bottom: var(--spacing-xl);
}

.login-footer {
  text-align: center;
  border-top: 1px solid var(--color-border);
  padding-top: var(--spacing-xl);
}

.login-footer p {
  margin: 0;
  font-size: var(--font-size-md);
}

.btn-login {
  background-color: var(--color-primary) !important;
  height: 44px;
  line-height: 44px;
  font-size: var(--font-size-lg);
  border-radius: var(--radius-md);
}

.btn-login:hover {
  background-color: var(--color-primary-hover) !important;
}

.input-field input:focus + label {
  color: var(--color-primary) !important;
}

.input-field input:focus {
  border-bottom: 1px solid var(--color-primary) !important;
  box-shadow: 0 1px 0 0 var(--color-primary) !important;
}

.input-field .prefix.active {
  color: var(--color-primary) !important;
}
</style> 