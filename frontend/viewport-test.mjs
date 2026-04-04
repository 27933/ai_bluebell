import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  
  // 测试不同的视口宽度
  const viewports = [
    { width: 1280, name: '1280' },
    { width: 1200, name: '1200' },
    { width: 1024, name: '1024' },
    { width: 768, name: '768' }
  ];

  try {
    for (const viewport of viewports) {
      const page = await browser.newPage();
      await page.setViewportSize({ width: viewport.width, height: 800 });
      await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
      await page.waitForTimeout(1500);
      
      // 获取container的实际宽度
      const containerWidth = await page.evaluate(() => {
        const container = document.querySelector('.container');
        if (!container) return 'N/A';
        return window.getComputedStyle(container).width;
      });

      await page.screenshot({ path: `/tmp/viewport-${viewport.name}.png`, fullPage: false });
      console.log(`✅ ${viewport.width}px: container宽度 = ${containerWidth}`);
      await page.close();
    }
  } finally {
    await browser.close();
  }
})();
