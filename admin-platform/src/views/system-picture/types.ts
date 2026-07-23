export interface PictureCategory {
  category_id: number;
  name: string;
}

export interface PictureFile {
  category_id: number;
  image: string;
  name: string;
  selected?: boolean;
}
