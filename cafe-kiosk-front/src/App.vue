<template>
  <div id="app">
    <router-view />
  </div>
</template>

<style>
:root {
  /* 북유럽 스타일 주요 색상 */
  --nordic-dark-grey: #2D3436;
  --nordic-light-grey: #F7F7F7;
  --nordic-pastel-blue: #78A1BB;
  --nordic-light-blue: #D6E2E9;
  --nordic-beige-brown: #BEB7A4;
  --nordic-grey-blue-light: #E2E6E8;
  --nordic-border-color: #D8DDE0;
  
  /* 기능별 색상 적용 */
  --text-primary: var(--nordic-dark-grey);
  --background-primary: var(--nordic-light-grey);
  --background-secondary: var(--nordic-grey-blue-light);
  --accent-color: var(--nordic-pastel-blue);
  --border-color: var(--nordic-border-color);
  --button-primary: var(--nordic-pastel-blue);
  --button-text: var(--nordic-light-grey);
  --card-background: white;

  --el-color-primary: var(--nordic-pastel-blue); /* 원하는 색상 코드로 변경 */
}
</style>