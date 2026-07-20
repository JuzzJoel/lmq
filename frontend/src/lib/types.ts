export interface Link {
  id: number;
  token: string;
  long_url: string;
  created_at: string;
  expires_at?: string;
  click_count: number;
  has_password?: boolean;
}

export interface LinkListResponse {
  links: Link[];
  total: number;
}

export interface ClickEvent {
  id: number;
  link_id: number;
  clicked_at: string;
  ip_address?: string;
  country: string;
  country_code: string;
  region: string;
  city: string;
  user_agent?: string;
  browser: string;
  os: string;
  is_mobile: boolean;
  referer?: string;
}

export interface ShortenRequest {
  url: string;
  custom_token?: string;
}

export interface ShortenResponse {
  token: string;
  short_url: string;
  long_url: string;
  created_at: string;
  has_password?: boolean;
  expires_at?: string;
}

export interface BulkShortenResponse {
  results: ShortenResponse[];
}

export interface APIResponse<T> {
  data: T;
  error: string | null;
  mock?: boolean;
}

export interface LinkAnalytics {
  token: string;
  long_url: string;
  total_clicks: number;
  clicks_by_day: DayCount[];
  countries: CountryCount[];
  country_groups: CountryCount[];
  regions: RegionCount[];
  cities: CityCount[];
  browsers: BrowserCount[];
  recent_clicks: ClickEvent[];
}

export interface DayCount {
  date: string;
  count: number;
}

export interface CountryCount {
  country: string;
  count: number;
}

export interface BrowserCount {
  browser: string;
  count: number;
}

export interface RegionCount {
  region: string;
  count: number;
}

export interface CityCount {
  city: string;
  count: number;
}
