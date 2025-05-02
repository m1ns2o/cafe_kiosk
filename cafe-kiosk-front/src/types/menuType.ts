export interface MenuItem {
  id: number;
  category_id: number;
  name: string;
  price: number;
  image_url: string;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}
