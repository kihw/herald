export const parseNumber = (v: any): number | undefined => {
  if (typeof v === 'number') return v;
  if (typeof v === 'string') {
    const n = parseFloat(v);
    return Number.isFinite(n) ? n : undefined;
  }
  return undefined;
};

export const inferBoolean = (v: any): boolean | null => {
  if (typeof v === 'boolean') return v;
  if (typeof v === 'string') {
    const lower = v.toLowerCase().trim();
    if (lower === 'true' || lower === '1' || lower === 'yes') return true;
    if (lower === 'false' || lower === '0' || lower === 'no') return false;
  }
  return null;
};