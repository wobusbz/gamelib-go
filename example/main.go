package main

import (
	"fmt"
	"gamelib-go/vector"
	"gamelib-go/zskip"
)

func ExampleVector() {
	v1 := vector.Vector2[float32]{X: 1, Y: 0} // →
	v2 := vector.Vector2[float32]{X: 0, Y: 1} // ↑
	v3 := vector.Vector2[float32]{X: 1, Y: 0}
	v4 := vector.Vector2[float32]{X: 1, Y: 1} // ↗
	v5 := vector.Vector2[float32]{X: 0, Y: 0}

	fmt.Printf("90° 夹角: %.1f°\n", vector.Angle(v1, v2))
	fmt.Printf("45° 夹角: %.1f°\n", vector.Angle(v3, v4))
	fmt.Printf("有符号 45°: %.1f°\n", vector.SignedAngle(v3, v4))
	fmt.Printf("有符号 -45°: %.1f°\n", vector.SignedAngle(v4, v3))
	fmt.Printf("零向量: %.1f°\n", vector.Angle(v1, v5))
	fmt.Printf("零向量: %.1f°\n", vector.SignedAngle(v1, v5))

	playerPos := vector.Vector2[float32]{X: 0, Y: 0}

	playerForward := vector.Vector2[float32]{X: 1, Y: 0}

	enemyPos := vector.Vector2[float32]{X: 5, Y: 5}

	fov := 90.0
	maxDist := 10.0

	if vector.InFOVDistance(playerPos, enemyPos, playerForward, fov, maxDist) {
		println("敌人在视野内")
	} else {
		println("敌人不在视野内")
	}
}

func ExampleSkipList() {
	skip := zskip.NewSkipList(10)
	skip.ZslInsert(1, "1")
	skip.ZslInsert(2, "2")
	skip.ZslInsert(3, "3")
	skip.TestPrint2()
	skip.DeleteFirst()
	skip.DeleteLast()
	skip.TestPrint2()
}

func ExampleSkipDict() {
	zskip.G_ZskDict.ZslSetEvictFront("1", 1, 3)
	zskip.G_ZskDict.ZslSetEvictFront("2", 2, 3)
	zskip.G_ZskDict.ZslSetEvictFront("3", 3, 3)
	zskip.G_ZskDict.ZslSetEvictFront("4", 4, 3)
	zskip.G_ZskDict.ZslSetEvictFront("5", 5, 3)
	zskip.G_ZskDict.TestPrint2()
}

func main() {
	ExampleSkipList()
	ExampleSkipDict()
}
