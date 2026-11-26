export function isPdfDocument() {
  return false;
}

export async function loadPdfDocument() {
  throw new Error("PDF viewer is disabled in MyPhotoX");
}

export function renderPdfPage() {
  return { promise: Promise.resolve(), cancel: () => {} };
}
