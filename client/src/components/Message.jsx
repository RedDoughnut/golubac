/*
  rounded - bool - if true all corners will be rounded, if false all except one will be rounded
  yours - bool - true if the message is by the user or false if by the contact
*/
export default function Message(props) {
  var yours = props.yours;
  var rounded = props.rounded;
  return (
    <div class="flex mr-2 ml-2 mt-0.5" classList={{
      "justify-end": yours
    }}>
      <div class="rounded-xl max-w-sm p-3 bg-(--primary) hover:bg-(--primary-hover) transition-colors duration-150 ease-out" classList={{
        "rounded-tr-none": !rounded && yours,
        "rounded-tl-none": !rounded && !yours
      }}>
        <span>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed tempor vel ligula et congue. Etiam finibus nulla vel dui sodales placerat. Sed vulputate ultricies ornare. Sed eu leo ultrices diam interdum aliquam ac a ipsum. Maecenas scelerisque sapien tortor. Vivamus venenatis ligula in quam consectetur tincidunt. Proin sed odio sed elit ornare rhoncus.</span>
      </div>
    </div>
  )
}
