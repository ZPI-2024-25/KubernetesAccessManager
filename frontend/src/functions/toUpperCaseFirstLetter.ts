export function capitalizeFirst<T extends string>(s: T){
    return s[0].toUpperCase() + s.slice(1)
}