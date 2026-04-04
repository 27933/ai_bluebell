import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page1 = await browser.newPage();
  const page2 = await browser.newPage();

  try {
    await page1.setViewportSize({ width: 1280, height: 800 });
    await page1.goto('file:///root/code/ai_bluebell/frontend_prototype.html', { waitUntil: 'networkidle' });
    await page1.screenshot({ path: '/tmp/proto-1280.png', fullPage: false });
    console.log('✅ 原型 1280px');

    await page2.setViewportSize({ width: 1280, height: 800 });
    await page2.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page2.waitForTimeout(1500);
    await page2.screenshot({ path: '/tmp/current-1280.png', fullPage: false });
    console.log('✅ 当前页 1280px');

  } finally {
    await browser.close();
  }
})();
