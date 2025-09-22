import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useMemo, useRef, useState } from "react";

export const Route = createFileRoute("/")({
	component: HomeComponent,
});

function HomeComponent() {
	const [apiStatus, setApiStatus] = useState<
		"loading" | "ok" | "error"
	>("loading");
	const [users, setUsers] = useState<
		{ id: string; username: string }[]
	>([]);
	const [myId, setMyId] = useState<string | null>(null);
	const [username, setUsername] = useState<string>("Guest");
	const [isConnected, setIsConnected] = useState<boolean>(false);
	const [message, setMessage] = useState<string>("");
	const [messages, setMessages] = useState<
		{
			username: string;
			userId?: string;
			content: string;
			type: string;
			timestamp?: string;
		}[]
	>([]);
	const wsRef = useRef<WebSocket | null>(null);

	const backendHttp = useMemo(
		() =>
			import.meta.env.VITE_BACKEND_URL?.replace(/\/$/, "") ||
			"http://localhost:8080",
		[]
	);
	const backendWs = useMemo(
		() => backendHttp.replace(/^http/, "ws"),
		[backendHttp]
	);

	// Health check
	useEffect(() => {
		let cancelled = false;
		fetch(`${backendHttp}/healthz`)
			.then((r) =>
				r.ok ? r.json() : Promise.reject(new Error("bad status"))
			)
			.then(() => !cancelled && setApiStatus("ok"))
			.catch(() => !cancelled && setApiStatus("error"));
		return () => {
			cancelled = true;
		};
	}, [backendHttp]);

	// Poll users when connected
	useEffect(() => {
		if (!isConnected) return;
		let cancelled = false;
		const load = () => {
			fetch(`${backendHttp}/api/users`)
				.then((r) =>
					r.ok ? r.json() : Promise.reject(new Error("bad status"))
				)
				.then((data) => {
					if (cancelled) return;
					if (Array.isArray(data?.usersDetailed)) {
						setUsers(data.usersDetailed);
					} else if (Array.isArray(data?.users)) {
						setUsers(
							data.users.map((u: string, idx: number) => ({
								id: String(idx + 1),
								username: u,
							}))
						);
					}
				})
				.catch(() => {});
		};
		load();
		const id = window.setInterval(load, 5000);
		return () => {
			cancelled = true;
			window.clearInterval(id);
		};
	}, [backendHttp, isConnected]);

	function connect() {
		if (wsRef.current || isConnected) return;
		const url = `${backendWs}/ws?username=${encodeURIComponent(username || "Guest")}`;
		const ws = new WebSocket(url);
		wsRef.current = ws;

		ws.onopen = () => {
			setIsConnected(true);
			pushMessage({
				username: "System",
				content: "Connected",
				type: "join",
			});
		};
		ws.onmessage = (ev) => {
			try {
				const data = JSON.parse(ev.data);
				if (data?.type === "welcome" && data?.userId) {
					setMyId(data.userId);
				}
				pushMessage({
					username: data.username ?? "Unknown",
					userId: data.userId,
					content: data.content ?? "",
					type: data.type ?? "chat",
					timestamp: data.timestamp,
				});
			} catch {
				pushMessage({
					username: "System",
					content: "Invalid message",
					type: "error",
				});
			}
		};
		ws.onclose = () => {
			setIsConnected(false);
			wsRef.current = null;
			pushMessage({
				username: "System",
				content: "Disconnected",
				type: "leave",
			});
		};
		ws.onerror = () => {
			pushMessage({
				username: "System",
				content: "WebSocket error",
				type: "error",
			});
		};
	}

	function disconnect() {
		wsRef.current?.close();
	}

	function send() {
		if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN)
			return;
		if (!message.trim()) return;
		const payload = {
			type: "chat",
			content: message.trim(),
			username,
		};
		wsRef.current.send(JSON.stringify(payload));
		setMessage("");
	}

	function pushMessage(m: {
		username: string;
		userId?: string;
		content: string;
		type: string;
		timestamp?: string;
	}) {
		setMessages((prev) => [...prev, m]);
	}

	return (
		<div className="mx-auto flex min-h-dvh max-w-5xl flex-col gap-4 px-4 py-6">
			<header className="flex items-center justify-between">
				<h1 className="text-lg font-semibold tracking-tight">
					Golang Realtime Chat
				</h1>
				<div className="text-sm">
					<span
						className={
							apiStatus === "ok" ?
								"rounded-full bg-green-100 px-2 py-1 text-green-700"
							: apiStatus === "error" ?
								"rounded-full bg-red-100 px-2 py-1 text-red-700"
							:	"rounded-full bg-gray-100 px-2 py-1 text-gray-600"
						}>
						{apiStatus === "loading" ?
							"Checking API..."
						:	apiStatus.toUpperCase()}
					</span>
				</div>
			</header>

			<section className="grid gap-4 md:grid-cols-3">
				<div className="rounded-xl border bg-white p-4 shadow-sm md:col-span-2">
					<div className="mb-3 flex items-center gap-2">
						<input
							type="text"
							value={username}
							onChange={(e) => setUsername(e.target.value)}
							placeholder="Username"
							className="w-48 rounded-lg border px-3 py-2 outline-none focus:ring"
						/>
						<button
							onClick={connect}
							disabled={isConnected}
							className="rounded-lg bg-black px-3 py-2 text-white transition disabled:opacity-50">
							Connect
						</button>
						<button
							onClick={disconnect}
							disabled={!isConnected}
							className="rounded-lg bg-neutral-200 px-3 py-2 text-neutral-900 transition hover:bg-neutral-300 disabled:opacity-50">
							Disconnect
						</button>
					</div>

					<div className="h-[480px] overflow-y-auto rounded-lg border bg-neutral-50 p-3">
						{messages.length === 0 ?
							<div className="flex h-full items-center justify-center text-sm text-neutral-500">
								No messages yet. Connect and say hello.
							</div>
						:	<ul className="space-y-2 text-sm">
								{messages.map((m, i) => {
									const isMine =
										m.userId && myId && m.userId === myId;
									const baseColor =
										m.type === "error" ? "text-red-700"
										: m.type === "join" ? "text-green-700"
										: m.type === "leave" ? "text-neutral-600"
										: "text-neutral-900";
									const bubble =
										isMine ?
											"bg-blue-600 text-white"
										:	"bg-white text-neutral-900 border";
									return (
										<li
											key={i}
											className={`flex ${isMine ? "justify-end" : "justify-start"}`}>
											<div
												className={`max-w-[75%] rounded-lg px-3 py-2 shadow-sm ${bubble}`}>
												<div
													className={`text-xs ${isMine ? "text-blue-100" : "text-neutral-500"}`}>
													{m.username} {myId}
												</div>
												<div className={baseColor}>{m.content}</div>
											</div>
										</li>
									);
								})}
							</ul>
						}
					</div>

					<div className="mt-3 flex items-center gap-2">
						<input
							type="text"
							value={message}
							onChange={(e) => setMessage(e.target.value)}
							placeholder="Type a message..."
							onKeyDown={(e) => {
								if (e.key === "Enter") send();
							}}
							className="flex-1 rounded-lg border bg-white px-3 py-2 outline-none focus:ring text-black"
						/>
						<button
							onClick={send}
							disabled={!isConnected || !message.trim()}
							className="rounded-lg bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:opacity-50">
							Send
						</button>
					</div>
				</div>

				<div className="rounded-xl border bg-white p-4 shadow-sm">
					<h3 className="mb-3 text-sm font-medium text-neutral-700">
						Connected Users
					</h3>
					<div className="min-h-10 rounded-lg border bg-neutral-50 p-3 text-sm">
						{users.length === 0 ?
							<div className="text-neutral-500">None</div>
						:	<div className="space-y-1">
								{users.map((u) => (
									<div
										key={u.id}
										className="flex items-center justify-between rounded bg-white px-2 py-1 shadow-sm text-black">
										<span>{u.username}</span>
										<span className="text-xs text-neutral-500">
											{u.id}
										</span>
									</div>
								))}
							</div>
						}
					</div>
				</div>
			</section>
		</div>
	);
}
