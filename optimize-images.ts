import sharp from "sharp";
import { readdirSync, mkdirSync, existsSync } from "node:fs";
import { join } from "node:path";

const picturesDir = join(import.meta.dir, "pictures");
const originalsDir = join(picturesDir, "originals");

// Create originals directory if it doesn't exist
if (!existsSync(originalsDir)) {
  mkdirSync(originalsDir, { recursive: true });
}

// Find all image files
const files = readdirSync(picturesDir).filter((f) =>
  f.match(/\.(jpg|jpeg|png)$/i)
);

if (files.length === 0) {
  console.log("No images found to optimize");
  process.exit(0);
}

console.log(`Found ${files.length} image(s) to optimize\n`);

for (const file of files) {
  const inputPath = join(picturesDir, file);
  const outputPath = join(
    picturesDir,
    file.replace(/\.(jpg|jpeg|png)$/i, ".webp")
  );

  // Skip if WebP already exists
  if (existsSync(outputPath)) {
    console.log(`⏭️  Skipping ${file} (WebP already exists)`);
    continue;
  }

  try {
    // Get original file size
    const originalStats = await Bun.file(inputPath).stat();
    const originalSizeMB = (originalStats.size / 1024 / 1024).toFixed(2);

    // Optimize and convert to WebP
    await sharp(inputPath)
      .resize(1200, null, {
        withoutEnlargement: true,
        fit: "inside",
      })
      .webp({
        quality: 85,
        effort: 6, // Higher effort = better compression (0-6)
      })
      .toFile(outputPath);

    // Get optimized file size
    const optimizedStats = await Bun.file(outputPath).stat();
    const optimizedSizeMB = (optimizedStats.size / 1024 / 1024).toFixed(2);
    const reduction = (
      ((originalStats.size - optimizedStats.size) / originalStats.size) *
      100
    ).toFixed(1);

    console.log(
      `✅ ${file} → ${file.replace(/\.(jpg|jpeg|png)$/i, ".webp")}`
    );
    console.log(
      `   ${originalSizeMB}MB → ${optimizedSizeMB}MB (${reduction}% reduction)\n`
    );

    // Move original to originals folder
    await Bun.write(
      join(originalsDir, file),
      await Bun.file(inputPath).arrayBuffer()
    );
    console.log(`📦 Moved original to pictures/originals/${file}\n`);
  } catch (error) {
    console.error(`❌ Error processing ${file}:`, error);
  }
}

console.log("✨ Image optimization complete!");
console.log(
  "\n💡 Don't forget to update your markdown files to use .webp extensions"
);
