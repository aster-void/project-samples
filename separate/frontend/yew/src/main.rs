use yew::prelude::*;

#[function_component]
fn App() -> Html {
    let counter = use_state(|| 0);
    let add = {
        let counter = counter.clone();
        let value = *counter + 1;
        move |_| {
            counter.set(value);
        }
    };
    let sub = {
        let counter = counter.clone();
        let value = *counter - 1;
        move |_| {
            counter.set(value);
        }
    };

    html! {
        <div>
            <button onclick={add}>{ "+1" }</button>
            <button onclick={sub}>{ "-1" }</button>
            <p>{ *counter }</p>
        </div>
    }
}

fn main() {
    yew::Renderer::<App>::new().render();
}
