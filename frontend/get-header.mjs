import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  try {
    await page.setViewportSize({ width: 1280, height: 400 });
    await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page.waitForTimeout(1500);
    await page.screenshot({ path: '/tmp/header-only.png', fullPage: false });
    console.log('✅ 头部截图已保存: /tmp/header-only.png');
  } catch (error) {
    console.error('❌ 出错:', error.message);
  } finally {
    await browser.close();
  }
})();
