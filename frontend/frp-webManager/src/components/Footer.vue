<template>
  	<footer class="footer">
    	<div class="footer-content">
			<div class="left">
            	<p class="timestamp">当前时间:{{ timestamp }}</p>
        	</div>
    	</div>
  	</footer>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const timestamp = ref('')

const fetchTimestamp = async () => {
    const response = await fetch('http://localhost:8080/api/time')
    const data = await response.json()
	timestamp.value = data.timestamp
}

let timer = null

onMounted(() => {
  fetchTimestamp()
  timer = setInterval(fetchTimestamp, 1000)//每秒
})

//组件卸载时清除定时器
onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})

</script>

<style scoped>
/*样式*/
.footer {
  background-color: #eceaea;
  padding: 0.05rem 0;
  position: fixed;
  bottom: 0;
  left:0;
  width: 100vw;
  box-shadow: 0 -3px 4px rgba(0, 0, 0, 0.479);
  z-index: 1000; /*层级设置*/
  transition: all 0.3s ease;
}

.footer-content {
  margin: 0 auto;
  padding: 0.4rem; /*边距*/
  display: flex;
  justify-content: space-between;
  align-items: center;
  justify-content: space-between;
  padding-left: 1rem
}


/*左列表*/
.left {
  display: flex; /*使用弹性布局*/
  list-style: none; /*移除列表默认样式*/
  align-items: center;
  gap: 1rem; /*间距*/
  margin: 0; /*外边距*/
  padding: 0; /*内边距*/
  font-size: 0.7rem;
}

.left .timestamp {
  text-decoration: none; /*移除下划线*/
  color: #000000; /*字体颜色*/
  font-weight: 500; /*字体粗细*/
}

/*小屏幕
@media (max-width: 768px) {
  .footer-content {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
}*/
</style>

<!-- 全局暗色主题样式 -->
<style>
/* 暗色主题 */
.dark-theme .footer {
  background-color: #2d2d2d;
  box-shadow: 0 -2px 4px rgba(255, 255, 255, 0.1);
}

.dark-theme .left .timestamp {
  text-decoration: none; /*移除下划线*/
  color: #ffffff; /*字体颜色*/
  font-weight: 100; /*字体粗细*/
}
</style>