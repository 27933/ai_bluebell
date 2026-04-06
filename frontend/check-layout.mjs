import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  
  const viewports = [
    { width: 1280, name: '1280' },
    { width: 1600, name: '1600' },
    { width: 1920, name: '1920' }
  ];

  try {
    for (const viewport of viewports) {
      const page = await browser.newPage();
      await page.setViewportSize({ width: viewport.width, height: 800 });
      await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
      await page.waitForTimeout(1500);
      
      const containerWidth = await page.evaluate(() => {
        const container = document.querySelector('.container');
        const containerStyle = window.getComputedStyle(container);
        const containerBox = container.getBoundingClientRect();
        return {
          width: containerStyle.width,
          marginLeft: containerStyle.marginLeft,
          marginRight: containerStyle.marginRight,
          paddingLeft: containerStyle.paddingLeft,
          paddingRight: containerStyle.paddingRight,
          actualWidth: containerBox.width
        };
      });

      console.log(`\n${viewport.width}px 视口:`);
      console.log(`  Container 宽度: ${containerWidth.width}`);
      console.log(`  Padding L/R: ${containerWidth.paddingLeft} / ${containerWidth.paddingRight}`);
      console.log(`  实际宽度: ${containerWidth.actualWidth}px`);
      
      await page.screenshot({ path: `/tmp/layout-${viewport.name}.png`, fullPage: false });
      await page.close();
    }
  } finally {
    await browser.close();
  }
})();
