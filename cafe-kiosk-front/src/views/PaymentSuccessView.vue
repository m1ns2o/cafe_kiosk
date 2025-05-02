<template>
  <div class="payment-success-container">
    <div class="payment-success">
      <div class="payment-info">
        <div class="success-status">
          <div class="status-icon">
            <span class="material-icon">check_circle</span>
          </div>
          <div class="status-message">결제가 성공했습니다!</div>
          <div class="status-description">잠시만 기다려주세요. 주문 페이지로 이동합니다.</div>
          <div class="countdown">{{ countdown }}초 후 이동</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const countdown = ref(3);

onMounted(() => {
  // 3초 후 주문 페이지로 자동 이동
  const timer = setInterval(() => {
    countdown.value--;
    if (countdown.value <= 0) {
      clearInterval(timer);
      router.push({ name: 'OrderView' });
    }
  }, 1000);
});
</script>

<style scoped>
.payment-success-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: var(--background-primary, #f5f5f5);
  padding: 20px;
}

.payment-success {
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.payment-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.success-status {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 15px;
  border-radius: 8px;
  background-color: #e6f7e6;
  margin-top: 10px;
  text-align: center;
}

.status-icon {
  font-size: 2rem;
  margin-bottom: 20px;
}

.status-icon .material-icon {
  font-family: 'Material Icons';
  font-weight: normal;
  font-style: normal;
  font-size: 80px;
  line-height: 1;
  letter-spacing: normal;
  text-transform: none;
  display: inline-block;
  white-space: nowrap;
  word-wrap: normal;
  direction: ltr;
  font-feature-settings: 'liga';
  -webkit-font-feature-settings: 'liga';
  -webkit-font-smoothing: antialiased;
  color: #4caf50;
}

.status-message {
  font-size: 1.8rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 15px;
}

.status-description {
  font-size: 1.2rem;
  color: #666;
  margin-bottom: 20px;
}

.countdown {
  font-size: 1.4rem;
  font-weight: 500;
  color: var(--button-primary, #4caf50);
  background-color: rgba(76, 175, 80, 0.1);
  padding: 10px 20px;
  border-radius: 30px;
  margin-top: 10px;
}

@media (max-width: 768px) {
  .payment-success {
    padding: 15px;
  }
  
  .status-icon .material-icon {
    font-size: 60px;
  }
  
  .status-message {
    font-size: 1.5rem;
  }
  
  .status-description {
    font-size: 1rem;
  }
  
  .countdown {
    font-size: 1.2rem;
  }
}

@media (max-width: 480px) {
  .payment-success-container {
    padding: 10px;
  }
  
  .success-status {
    padding: 30px 10px;
  }
  
  .status-icon .material-icon {
    font-size: 50px;
  }
  
  .status-message {
    font-size: 1.3rem;
  }
}
</style>