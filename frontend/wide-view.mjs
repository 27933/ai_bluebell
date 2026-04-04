import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  try {
    await page.setViewportSize({ width: 1600, height: 900 });
    await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page.waitForTimeout(2000);
    await page.screenshot({ path: '/tmp/wide-view.png', fullPage: true });
    console.log('✅ 宽视图已保存: /tmp/wide-view.png');
  } finally {
    await browser.close();
  }
})();
