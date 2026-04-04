import { chromium } from '@playwright/test';

(async () => {
  console.log('对比原型和现在的前端...\n');

  const browser = await chromium.launch({ headless: true });
  const page1 = await browser.newPage();
  const page2 = await browser.newPage();

  try {
    // 打开原型
    console.log('📄 正在加载原型 HTML...');
    await page1.goto('file:///root/code/ai_bluebell/frontend_prototype.html', { waitUntil: 'networkidle' });
    await page1.screenshot({ path: '/tmp/prototype.png', fullPage: false });
    console.log('✅ 原型截图已保存: /tmp/prototype.png');

    // 打开现在的前端
    console.log('\n🔄 正在加载当前前端...');
    await page2.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page2.waitForTimeout(2000);
    await page2.screenshot({ path: '/tmp/current-frontend.png', fullPage: false });
    console.log('✅ 当前前端截图已保存: /tmp/current-frontend.png');

  } catch (error) {
    console.error('❌ 出错:', error.message);
  } finally {
    await browser.close();
  }
})();
