import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page1 = await browser.newPage();
  const page2 = await browser.newPage();

  try {
    // 统一使用相同的视口大小
    const size = { width: 1280, height: 800 };
    
    await page1.setViewportSize(size);
    await page1.goto('file:///root/code/ai_bluebell/frontend_prototype.html', { waitUntil: 'networkidle' });
    await page1.screenshot({ path: '/tmp/proto-compare.png', fullPage: false });
    console.log('✅ 原型截图');

    await page2.setViewportSize(size);
    await page2.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page2.waitForTimeout(2000);
    await page2.screenshot({ path: '/tmp/current-compare.png', fullPage: false });
    console.log('✅ 当前页面截图');

  } finally {
    await browser.close();
  }
})();
