<template>
  <div class="notification-container">
    <transition-group name="notification">
      <div
        v-for="toast in notifications"
        :key="toast.id"
        class="notification"
        :class="`notification-${toast.type}`"
      >
        <i class="material-icons notification-icon">{{ iconMap[toast.type] }}</i>
        <span class="notification-message">{{ toast.message }}</span>
      </div>
    </transition-group>
  </div>
</template>

<script>
export default {
  name: 'NotificationContainer',
  data() {
    return {
      notifications: [],
      iconMap: {
        success: 'check_circle',
        error: 'error',
        info: 'info',
        warning: 'warning'
      }
    };
  },
  methods: {
    show(message, type = 'info', duration = 2500) {
      const id = Date.now() + Math.random();
      this.notifications.push({ id, message, type });
      setTimeout(() => {
        this.notifications = this.notifications.filter(n => n.id !== id);
      }, duration);
    }
  }
};
</script>

<style>
.notification-container {
  position: fixed;
  top: 80px;
  right: 24px;
  z-index: 11000;
  display: flex;
  flex-direction: column;
  gap: 10px;
  pointer-events: none;
}

.notification {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  background: #FFFFFF;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12), 0 0 0 1px rgba(0, 0, 0, 0.04);
  min-width: 240px;
  max-width: 380px;
  pointer-events: auto;
  font-family: Inter, "PingFang SC", "Microsoft YaHei", sans-serif;
}

.notification-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.notification-message {
  font-size: 14px;
  color: #0F172A;
  line-height: 1.4;
}

.notification-success .notification-icon { color: #16A34A; }
.notification-success { border-left: 3px solid #16A34A; }

.notification-error .notification-icon { color: #EF4444; }
.notification-error { border-left: 3px solid #EF4444; }

.notification-info .notification-icon { color: #2563EB; }
.notification-info { border-left: 3px solid #2563EB; }

.notification-warning .notification-icon { color: #F59E0B; }
.notification-warning { border-left: 3px solid #F59E0B; }

/* 动画 */
.notification-enter-active {
  transition: all 0.3s ease;
}
.notification-leave-active {
  transition: all 0.25s ease;
}
.notification-enter-from {
  opacity: 0;
  transform: translateX(30px);
}
.notification-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
