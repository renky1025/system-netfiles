import { test, expect } from '@playwright/test';

const TEST_USER = `testuser_${Date.now()}`;
const TEST_PASS = 'password123';

test.describe('File Operations', () => {

    test('should register and login', async ({ page }) => {
        await page.goto('/register');
        await page.getByPlaceholder('Username').fill(TEST_USER);
        await page.getByPlaceholder('Email').fill('test@example.com');
        await page.getByPlaceholder('Password').fill(TEST_PASS);

        await page.click('button:has-text("Register")');

        // Should redirect to login
        await expect(page).toHaveURL('/login');

        // Login
        await page.getByPlaceholder('Username').fill(TEST_USER);
        await page.getByPlaceholder('Password').fill(TEST_PASS);
        await page.click('button:has-text("Login")');

        // Should redirect to home
        await expect(page).toHaveURL('/');
    });

    test('should upload file and update list', async ({ page }) => {
        // Login first
        await page.goto('/login');
        await page.getByPlaceholder('Username').fill(TEST_USER);
        await page.getByPlaceholder('Password').fill(TEST_PASS);
        await page.click('button:has-text("Login")');
        await expect(page).toHaveURL('/');

        // Create a dummy file
        // @ts-ignore
        const buffer = Buffer.from('Hello World');

        // Upload
        const fileInput = page.locator('input[type="file"]');
        await fileInput.setInputFiles({
            name: 'test_root.txt',
            mimeType: 'text/plain',
            buffer: buffer,
        });

        // Verify toast or list update
        await expect(page.locator('.file-name-cell', { hasText: 'test_root.txt' })).toBeVisible();
    });

    test('should create folder and verify tree update', async ({ page }) => {
        await page.goto('/login');
        await page.getByPlaceholder('Username').fill(TEST_USER);
        await page.getByPlaceholder('Password').fill(TEST_PASS);
        await page.click('button:has-text("Login")');

        await page.click('button:has-text("New Folder")');
        // FolderDialog has input with placeholder "Folder Name" or we can use selector
        // Let's use a more generic selector if unsure
        await page.locator('.el-dialog input').fill('Docs');

        // Click Create in the dialog
        await page.locator('.el-dialog__footer button:has-text("Create")').click();

        // Verify in list
        await expect(page.locator('.file-name-cell', { hasText: 'Docs' })).toBeVisible();

        // Verify in tree (sidebar)
        await expect(page.locator('.el-tree-node__label', { hasText: 'Docs' })).toBeVisible();
    });

    test('should upload to subfolder and verify display', async ({ page }) => {
        await page.goto('/login');
        await page.getByPlaceholder('Username').fill(TEST_USER);
        await page.getByPlaceholder('Password').fill(TEST_PASS);
        await page.click('button:has-text("Login")');

        // Enter folder
        await page.dblclick('.file-name-cell:has-text("Docs")');

        // Upload file
        // @ts-ignore
        const buffer = Buffer.from('Sub Content');
        const fileInput = page.locator('input[type="file"]');
        await fileInput.setInputFiles({
            name: 'sub_file.txt',
            mimeType: 'text/plain',
            buffer: buffer,
        });

        // Verify
        await expect(page.locator('.file-name-cell', { hasText: 'sub_file.txt' })).toBeVisible();

        // Go back
        await page.click('.el-breadcrumb__item:has-text("Home")');

        // Verify we are back in root
        await expect(page.locator('.file-name-cell', { hasText: 'Docs' })).toBeVisible();
        await expect(page.locator('.file-name-cell', { hasText: 'sub_file.txt' })).not.toBeVisible();
    });
});
