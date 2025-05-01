
<template>
  <div class="page-component">
    <h1>Home Page</h1>
    <p>This is the home page content</p>
  </div>
</template>

<style scoped>
.page-component {
  height: 100%;
  width: 100%;
  box-sizing: border-box;
  overflow: auto; /* 필요한 경우에만 내부 스크롤 허용 */
}

/* 글로벌 CSS 정의 (필요한 경우) */
* {
  box-sizing: border-box;
}

body, html {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
}

#app {
  font-family: Arial, sans-serif;
  height: 100%;
  width: 100%;
}

/* router-view가 전체 화면을 차지하도록 설정 */
.router-view-container {
  height: 100%;
  width: 100%;
}
</style>