import type { APIResponse, Link, LinkAnalytics, ShortenRequest, BulkShortenResponse } from './types';

const API_BASE = '/api/v1';

/** Creates a shortened URL */
export async function shortenUrl(request: ShortenRequest): Promise<APIResponse<BulkShortenResponse>> {
  try {
    const response = await fetch(`${API_BASE}/shorten`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request)
    });
    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        return { data: {} as BulkShortenResponse, error: errorData.error || 'Failed to shorten URL' };
    }
    return response.json() as Promise<APIResponse<BulkShortenResponse>>;
  } catch (error) {
    return { data: {} as BulkShortenResponse, error: 'Network error. Please check your connection.' };
  }
}

/** Fetches analytics for a specific token */
export async function getAnalytics(token: string): Promise<APIResponse<LinkAnalytics>> {
  try {
    const response = await fetch(`${API_BASE}/analytics?token=${encodeURIComponent(token)}`, {
      headers: {
        'X-Admin-Token': typeof sessionStorage !== 'undefined' ? sessionStorage.getItem('admin_token') || '' : ''
      }
    });
    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        return { data: {} as LinkAnalytics, error: errorData.error || 'Failed to fetch analytics' };
    }
    return response.json() as Promise<APIResponse<LinkAnalytics>>;
  } catch (error) {
    return { data: {} as LinkAnalytics, error: 'Network error.' };
  }
}

/** Lists paginated links with summary data */
export async function getLinks(page: number = 1, limit: number = 20, search: string = ''): Promise<APIResponse<import('./types').LinkListResponse>> {
  try {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString()
    });
    if (search) params.append('search', search);

    const response = await fetch(`${API_BASE}/analytics/links?${params.toString()}`, {
      headers: {
        'X-Admin-Token': typeof sessionStorage !== 'undefined' ? sessionStorage.getItem('admin_token') || '' : ''
      }
    });
    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        return { data: { links: [], total: 0 }, error: errorData.error || 'Failed to fetch links' };
    }
    return response.json() as Promise<APIResponse<import('./types').LinkListResponse>>;
  } catch (error) {
    return { data: { links: [], total: 0 }, error: 'Network error.' };
  }
}
