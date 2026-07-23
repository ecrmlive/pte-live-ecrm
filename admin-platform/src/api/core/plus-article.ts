export async function getArticleCategoryListApi() {
  return { list: [] as Array<{ category_id: number; name: string }> };
}
