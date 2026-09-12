export default function Contact(props) {
  return (
      <div class="p-3 flex items-center rounded-2xl hover:opacity-75 cursor-pointer">
        <img src="assets/pfp.jpg" class="rounded-full" width="40" />
        <div class="ml-2">
          <p>{props.name}</p>
          <p class="text-sm">Start the chat</p>
        </div>
      </div>
  )
}
