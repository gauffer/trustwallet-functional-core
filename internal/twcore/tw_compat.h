// Заголовки wallet-core рассчитывают на наличие Clang built-in макросов
// (__has_feature, _Nonnull и др.). Если их нет, то препроцессор падает.
// Шим снимает зависимость от конкретного компилятора.
// На бинарник не влияет, всё раскрывается в ноль или пусто.
#ifndef __has_feature
#define __has_feature(x) 0
#endif
#ifndef __has_extension
#define __has_extension(x) 0
#endif
#ifndef _Nonnull
#define _Nonnull
#endif
#ifndef _Nullable
#define _Nullable
#endif
