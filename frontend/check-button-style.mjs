import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  try {
    await page.setViewportSize({ width: 1600, height: 900 });
    await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page.waitForTimeout(2000);

    const button = await page.$('.btn-primary');
    if (button) {
      const bgColor = await button.evaluate(el => window.getComputedStyle(el).backgroundColor);
      const color = await button.evaluate(el => window.getComputedStyle(el).color);
      const opacity = await button.evaluate(el => window.getComputedStyle(el).opacity);
      const zIndex = await button.evaluate(el => window.getComputedStyle(el).zIndex);
      const overflow = await button.evaluate(el => window.getComputedStyle(el).overflow);
      
      console.log('按钮样式:');
      console.log('  backgroundColor:', bgColor);
      console.log('  color:', color);
      console.log('  opacity:', opacity);
      console.log('  z-index:', zIndex);
      console.log('  overflow:', overflow);
    }
  } finally {
    await browser.close();
  }
})();
