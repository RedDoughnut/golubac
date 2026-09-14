export default function Setting(props) {
  return (
    <button class="mx-1 p-2 hover:bg-(--primary-hover) rounded-md w-full cursor-pointer text-left">
      <p class="text-md font-semibold">{props.name}</p>
    </button>
  )
}
