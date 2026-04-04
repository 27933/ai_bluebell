import { chromium } from '@playwright/test';

(async () => {
  console.log('启动浏览器进行调试...\n');

  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  // 监听控制台消息
  page.on('console', msg => console.log(`[Console] ${msg.text()}`));

  // 监听请求
  page.on('request', request => {
    if (request.url().includes('/api/')) {
      console.log(`[Request] ${request.method()} ${request.url()}`);
    }
  });

  // 监听响应
  page.on('response', response => {
    if (response.url().includes('/api/')) {
      console.log(`[Response] ${response.status()} ${response.url()}`);
    }
  });

  try {
    await page.setViewportSize({ width: 1200, height: 800 });
    await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page.waitForTimeout(3000);

    // 检查DOM中是否有文章卡片
    const cardCount = await page.$$('.article-card').then(cards => cards.length);
    console.log(`\n📊 页面中的文章卡片数: ${cardCount}`);

    // 检查 empty-state
    const emptyState = await page.$('.empty-state');
    if (emptyState) {
      const emptyText = await emptyState.textContent();
      console.log(`空状态文本: "${emptyText?.trim()}"`);
    }

    // 尝试获取 articles ref 的值
    const articlesCount = await page.evaluate(() => {
      // 从 Vue 实例中获取数据
      const app = document.querySelector('#app').__vue_app__;
      // 这可能不工作，但尝试一下
      return 'see console for details';
    }).catch(() => 'unable to access Vue instance');

    console.log(`Vue data: ${articlesCount}`);

    // 截图保存
    await page.screenshot({ path: '/tmp/debug-home-page.png', fullPage: true });
    console.log('\n✅ 调试截图已保存: /tmp/debug-home-page.png');

  } catch (error) {
    console.error('❌ 出错:', error.message);
  } finally {
    await browser.close();
  }
})();
