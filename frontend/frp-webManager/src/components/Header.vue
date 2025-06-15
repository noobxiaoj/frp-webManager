<template>
  	<header class="header">
    	<nav class="nav-container">
      	<!-- Logo和主题切换的容器 -->
      		<div class="logo-section">
        		<!-- Logo -->
        		<div class="logo">FRPS-WEB</div>
        
        		<!-- 主题切换 -->

			
        		<div class="theme-toggle-out" @click="toggleTheme">
					<button class="theme-toggle">
						<span v-if="isDark">🌙</span>
          				<span v-else>🌞</span>
					</button>
        		</div>
      		</div>
      
      		<!-- 导航链接列表 -->
      		<ul class="nav-links">
        		<li><a href="/home">首页</a></li>
        		<li><a href="/about">控制台</a></li>
				<div class="user_out">
					<a href="/user" class = "user">
						<div class="gg-boy"></div>
					</a>
					<a href="#logout" class="logout">登出</a>
				</div>
      		</ul>
    	</nav>
  	</header>
</template>

<script setup>
import { ref, onMounted } from 'vue'

//
//false为浅色主题
const isDark = ref(false)

//切换主题
const toggleTheme = () => {
  isDark.value = !isDark.value
  
  if (isDark.value) {
    //添加dark-theme类到body
    document.body.classList.add('dark-theme')
    //保存到本地
    localStorage.setItem('theme', 'dark')
  } else {
    //移除dark-theme类
    document.body.classList.remove('dark-theme')
    //保存到本地
    localStorage.setItem('theme', 'light')
  }
}

//组件挂载时执行的函数
onMounted(() => {
  //获取保存的偏好
  const savedTheme = localStorage.getItem('theme')
  
  if (savedTheme === 'dark') {
    //设置isDark为true
    isDark.value = true
    document.body.classList.add('dark-theme')
  }
})
</script>

<style scoped>
.header {
  background-color: #eceaea;/*背景颜色*/
  box-shadow: 0 3px 4px rgba(0, 0, 0, 0.479); /*阴影*/
  position: fixed; /*始终在页面顶部*/
  width: 100vw; 
  left:0;
  top: 0;
  z-index: 1000; /*层级设置*/
  transition: all 0.3s ease; /*添加过渡动画*/
}

/*导航容器*/
.nav-container {
  margin: 0 auto; /*居中*/
  padding: 1rem; /*边距*/
  display: flex; /*使用弹性布局*/
  justify-content: space-between; /*两端对齐*/
  align-items: center; /*垂直居中对齐*/
}

/*Logo区域*/
.logo-section {
  display: flex; 
  align-items: center; 
  gap: 1rem; 
}

/*Logo字体*/
.logo {
  font-size: 1rem; /*大小*/
  font-weight: bold; /*加粗*/
  color: #ff0000; /*颜色*/
}

/*导航列表*/
.nav-links {
  display: flex; /*使用弹性布局*/
  list-style: none; /*移除列表默认样式*/
  align-items: center;
  gap: 2rem; /*间距*/
  margin: 0; /*外边距*/
  padding: 0; /*内边距*/
}

/*导航链接样式*/
.nav-links a {
  text-decoration: none; /*移除下划线*/
  color: #000000; /*字体颜色*/
  font-weight: 500; /*字体粗细*/
  transition: color 0.3s ease; /*颜色变化的过渡动画*/
}

/*悬停效果*/
.nav-links a:hover {
  color: #ff0000; /*悬停颜色*/
}

/*主题切换*/
.theme-toggle-out {
  background: none; /*无背景色*/
  border: 1px solid #333; /*边框样式*/
  box-shadow: 2px 3px 4px rgba(0, 0, 0, 0.671); /*阴影*/
  border-radius: 50px; /*圆形按钮*/
  padding:0px 25px;
  width: 0%;
  height: 30px;
  cursor: pointer; /*鼠标悬停时显示手型光标*/
  display: flex;
  align-items: center;
  justify-content: flex-end;
  transition: all 0.3s ease; /*过渡动画*/
  flex-shrink: 0; /*防止被压缩*/
}


.theme-toggle {
  background: none; /*无背景色*/
  border: 0px solid #333; /*边框样式*/
  box-shadow: 3px 0px 4px rgba(0, 0, 0, 0.671); /*阴影*/
  border-radius: 50%; /*圆形按钮*/
  width: 25px; /*宽度*/
  height: 25px; /*高度*/
  cursor: pointer; /*鼠标悬停时显示手型光标*/
  font-size: 1.3rem; /*内容大小*/
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease; /* 过渡动画 */
}

.user_out {
  background: none; /*无背景色*/
  border: 1px solid #333; /*边框样式*/
  border-radius: 50px; /*圆形按钮*/
  width: 100px; /*宽度*/
  height: 50px; /*高度*/
  display: flex; /*使用弹性布局*/
  align-items: center;
  justify-content: space-between;
  list-style: none; /*移除列表默认样式*/
  padding-right: 5px
}

.user {
  background: rgb(255, 95, 122); /*无背景色*/
  border: 1px solid #333; /*边框样式*/
  border-radius: 50px; /*圆形按钮*/
  width: 50px; /*宽度*/
  height: 50px; /*高度*/
  display: flex; /*使用弹性布局*/
  align-items: center;
  justify-content: center;
  list-style: none; /*移除列表默认样式*/
}

.gg-boy,
.gg-boy::after,
.gg-boy::before {
  display: block;
  box-sizing: border-box;
  border-radius: 42px;
}
.gg-boy {
  position: relative;
  width: 20px;
  height: 20px;
  transform: scale(2);
  overflow: hidden;
  box-shadow: inset 0 0 0 2px #000000;
}
.gg-boy::after,
.gg-boy::before {
  content: "";
  position: absolute;
  width: 2px;
  height: 2px;
  background: currentColor;
  box-shadow: 6px 0 0;
  left: 6px;
  top: 10px;
}
.gg-boy::after {
  width: 20px;
  height: 20px;
  top: -13px;
  right: -12px;
}

.user_out .logout {
  color: #00a6ac;
}

.user_out:hover .logout {
  color: #ff0000;
}

.user_out .avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/*小屏幕*/
/*
@media (max-width: 768px) {
  .nav-container {
    flex-direction: column;
    gap: 1rem;
	padding: 1rem;
  }
  .logo-section {
    align-self: flex-start;
  }
  .nav-links {
    gap: 1rem;
  }
}*/
</style>

<!-- 全局样式主题切换 -->
<style>
/*暗色主题样式*/
.dark-theme {
  background-color: #1a1a1a; /*背景*/
  color: #ffffff; /*文字*/
}

.dark-theme .header {
  background-color: #2d2d2d;
  box-shadow: 0 2px 4px rgba(255,255,255,0.1);
}

.dark-theme .logo {
  color: #00a6ac;
}

.dark-theme .nav-links a {
  color: #ffffff;
}

.dark-theme .nav-links a:hover {
  color: #00a6ac; /*悬停颜色*/
}

.dark-theme .theme-toggle-out {
  border: 1px solid #ffffff; /*边框样式*/
  justify-content: flex-start;
  box-shadow: -2px 3px 3px rgba(255, 255, 255, 0.5); /*阴影*/
}

.dark-theme .theme-toggle {
  border-color: #ffffff; /* 白色边框 */
  color: #ffffff; /* 白色图标 */
  box-shadow: -2px 3px 4px rgb(255, 255, 255); /*阴影*/
}

.dark-theme .user_out {
  border: 1px solid #ffffff; /*边框样式*/
}

.dark-theme .user {
  background: rgb(80, 17, 133);
  border: 1px solid #ffffff; /*边框样式*/
}

.dark-theme .gg-boy {
  box-shadow: inset 0 0 0 2px #ffffff; /* 夜间模式下boy的颜色为白色 */
}

.dark-theme .user_out .logout {
  color: #ff0000;
}

.dark-theme .user_out:hover .logout {
  color: #00a6ac;
}

</style>